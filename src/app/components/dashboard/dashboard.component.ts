import { AfterViewInit, Component, ElementRef, OnDestroy, QueryList, ViewChildren, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MapViewComponent } from './map-view.component/map-view.component';
import { ATCCComponent } from '../atcc/atcc.component';
import { DashboardDataService, MenuItem } from '../../services/dashboard-data.service';
import { AtccSection } from '../../services/atcc-data.service';

interface MpegtsPlayer {
  attachMediaElement(video: HTMLVideoElement): void;
  load(): void;
  play(): Promise<void> | void;
  unload(): void;
  detachMediaElement(): void;
  destroy(): void;
}

interface MpegtsFactory {
  isSupported(): boolean;
  createPlayer(mediaDataSource: { type: 'flv'; url: string; isLive: boolean }, config?: Record<string, unknown>): MpegtsPlayer;
}

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
  imports: [
    CommonModule,
    MapViewComponent,
    ATCCComponent,
  ],
})
export class DashboardComponent implements AfterViewInit, OnDestroy {
  @ViewChildren('flvVideo') private flvVideos?: QueryList<ElementRef<HTMLVideoElement>>;

  private readonly dashboardDataService = inject(DashboardDataService);
  private readonly dashboardData = this.dashboardDataService.getDashboardData();
  private readonly flvPlayers = new Map<HTMLVideoElement, MpegtsPlayer>();
  private readonly flvStartTimers = new Map<HTMLVideoElement, number>();
  private readonly bufferedPlaybackStarted = new WeakSet<HTMLVideoElement>();
  private mpegtsLoader: Promise<MpegtsFactory | null> | null = null;

  isLightTheme = false;
  isSidebarHidden = false;

  readonly notificationCount = this.dashboardData.header.notificationCount;
  readonly userName = this.dashboardData.header.userName;
  readonly clock = this.dashboardData.header.clock;
  readonly sideMenuItems = this.dashboardData.navigation.sideMenuItems;
  readonly topTabs = this.dashboardData.navigation.topTabs;
  readonly kpis = this.dashboardData.kpis;
  readonly mapLegend = this.dashboardData.mapLegend;
  readonly videoFeeds = this.dashboardData.videoFeeds;
  readonly alerts = this.dashboardData.alerts;
  readonly vmsMessages = this.dashboardData.vmsMessages;
  readonly weatherStats = this.dashboardData.weatherStats;
  readonly vsdsMetrics = this.dashboardData.vsdsMetrics;
  readonly incidents = this.dashboardData.incidents;
  readonly deviceSummary = this.dashboardData.deviceSummary;
  readonly deviceHealthMetrics = this.dashboardData.deviceHealthMetrics;
  readonly deviceStatus = this.dashboardData.deviceStatus;
  readonly activities = this.dashboardData.activities;

  activeBody = 'dashboard';
  activeAtccSection: AtccSection = 'visualization';
  expandedMenu: string | null = null;
  activeTopTab = 'Dashboard';

  ngAfterViewInit(): void {
    this.setupFlvPlayers();
    this.flvVideos?.changes.subscribe(() => this.setupFlvPlayers());
  }

  ngOnDestroy(): void {
    this.destroyFlvPlayers();
  }

  selectMenu(item: MenuItem): void {
    if (!item.children?.length) {
      if (item.route !== undefined) {
        this.activeBody = item.route;
      }
      this.activeTopTab = item.label;
      this.expandedMenu = null;
      return;
    }

    this.expandedMenu = this.expandedMenu === item.label ? null : item.label;
    this.activeTopTab = item.label;
  }

  selectLeaf(item: MenuItem): void {
    if (!item.route) {
      return;
    }

    this.activeBody = item.route;
    this.activeTopTab = 'ATCC';

    if (item.route.startsWith('atcc/')) {
      const section = item.route.split('/')[1] as AtccSection | undefined;
      this.activeAtccSection = section ?? 'visualization';
      this.expandedMenu = 'ATCC';
    }
  }

  hasActiveLeaf(item: MenuItem): boolean {
    if (item.route && item.route === this.activeBody) {
      return true;
    }

    return (item.children ?? []).some(child => this.hasActiveLeaf(child));
  }

  toggleSidebar(): void {
    this.isSidebarHidden = !this.isSidebarHidden;
  }

  toggleTheme(): void {
    this.isLightTheme = !this.isLightTheme;
  }

  toggleFullscreen(): void {
    const documentElement = document.documentElement;

    if (!document.fullscreenElement) {
      documentElement.requestFullscreen?.();
      return;
    }

    document.exitFullscreen?.();
  }

  sendPtzCommand(cameraId: string, command: 'left' | 'right' | 'up' | 'down' | 'zoomIn' | 'zoomOut'): void {
    void fetch(`/api/ptz/${encodeURIComponent(cameraId)}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ command }),
    }).catch(() => undefined);
  }

  applyLiveBuffer(event: Event, bufferSeconds = 2): void {
    const video = event.target as HTMLVideoElement;
    const liveEdge = this.getLiveEdge(video);

    if (liveEdge !== null) {
      video.currentTime = Math.max(0, liveEdge - bufferSeconds);
    }

    this.playWhenBuffered(event, bufferSeconds);
  }

  playWhenBuffered(event: Event, bufferSeconds = 5): void {
    const video = event.target as HTMLVideoElement;
    const requiredBuffer = Math.max(1, bufferSeconds);

    if (this.bufferedPlaybackStarted.has(video) || this.getBufferedAhead(video) < requiredBuffer) {
      return;
    }

    this.bufferedPlaybackStarted.add(video);

    const timer = this.flvStartTimers.get(video);
    if (timer !== undefined) {
      window.clearTimeout(timer);
      this.flvStartTimers.delete(video);
    }

    if (video.paused) {
      void video.play().catch(() => undefined);
    }
  }

  maintainLiveBuffer(event: Event, bufferSeconds = 2): void {
    const video = event.target as HTMLVideoElement;
    const liveEdge = this.getLiveEdge(video);

    if (liveEdge === null) {
      return;
    }

    const targetTime = Math.max(0, liveEdge - bufferSeconds);
    const drift = targetTime - video.currentTime;

    if (drift > 3 || drift < -1) {
      video.currentTime = targetTime;
    }

    if (video.paused && this.getBufferedAhead(video) >= Math.max(1, bufferSeconds)) {
      void video.play().catch(() => undefined);
    }
  }

  private getLiveEdge(video: HTMLVideoElement): number | null {
    const seekable = video.seekable;

    if (!seekable.length) {
      return null;
    }

    return seekable.end(seekable.length - 1);
  }

  private getBufferedAhead(video: HTMLVideoElement): number {
    const buffered = video.buffered;

    for (let index = 0; index < buffered.length; index += 1) {
      if (buffered.start(index) <= video.currentTime && buffered.end(index) >= video.currentTime) {
        return buffered.end(index) - video.currentTime;
      }
    }

    return 0;
  }

  private setupFlvPlayers(): void {
    const videos = this.flvVideos?.toArray().map(item => item.nativeElement) ?? [];
    const activeVideos = new Set(videos);

    this.flvPlayers.forEach((player, video) => {
      if (!activeVideos.has(video)) {
        this.destroyFlvPlayer(video, player);
      }
    });

    videos.forEach(video => {
      if (!this.flvPlayers.has(video)) {
        void this.createFlvPlayer(video);
      }
    });
  }

  private async createFlvPlayer(video: HTMLVideoElement): Promise<void> {
    const streamUrl = video.dataset['streamUrl'];

    if (!streamUrl) {
      return;
    }

    const mpegts = await this.loadMpegts();

    if (!mpegts?.isSupported()) {
      return;
    }

    const bufferSeconds = Number(video.dataset['bufferSeconds'] ?? 2);
    const player = mpegts.createPlayer(
      { type: 'flv', url: streamUrl, isLive: true },
      {
        enableWorker: true,
        enableStashBuffer: true,
        stashInitialSize: 1024 * 1024,
        liveBufferLatencyChasing: true,
        liveBufferLatencyMaxLatency: bufferSeconds + 1,
        liveBufferLatencyMinRemain: Math.max(0.5, bufferSeconds - 0.5),
      },
    );

    player.attachMediaElement(video);
    player.load();
    this.flvPlayers.set(video, player);

    const timer = window.setTimeout(() => {
      if (!this.bufferedPlaybackStarted.has(video)) {
        this.bufferedPlaybackStarted.add(video);
        const playResult = player.play();
        if (playResult instanceof Promise) {
          void playResult.catch(() => undefined);
        }
      }
    }, Math.max(1, bufferSeconds) * 1000);
    this.flvStartTimers.set(video, timer);
  }

  private loadMpegts(): Promise<MpegtsFactory | null> {
    if (this.mpegtsLoader) {
      return this.mpegtsLoader;
    }

    const existing = (window as unknown as { mpegts?: MpegtsFactory }).mpegts;
    if (existing) {
      this.mpegtsLoader = Promise.resolve(existing);
      return this.mpegtsLoader;
    }

    this.mpegtsLoader = new Promise(resolve => {
      const script = document.createElement('script');
      script.src = 'https://cdn.jsdelivr.net/npm/mpegts.js@1.7.3/dist/mpegts.min.js';
      script.async = true;
      script.onload = () => resolve((window as unknown as { mpegts?: MpegtsFactory }).mpegts ?? null);
      script.onerror = () => resolve(null);
      document.head.appendChild(script);
    });

    return this.mpegtsLoader;
  }

  private destroyFlvPlayers(): void {
    this.flvPlayers.forEach((player, video) => this.destroyFlvPlayer(video, player));
  }

  private destroyFlvPlayer(video: HTMLVideoElement, player: MpegtsPlayer): void {
    const timer = this.flvStartTimers.get(video);

    if (timer !== undefined) {
      window.clearTimeout(timer);
      this.flvStartTimers.delete(video);
    }

    player.unload();
    player.detachMediaElement();
    player.destroy();
    this.flvPlayers.delete(video);
  }

}
