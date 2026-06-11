import { AfterViewInit, Component, ElementRef, OnDestroy, OnInit, QueryList, ViewChildren, inject } from '@angular/core';
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
  on?(event: string, handler: (...args: unknown[]) => void): void;
  off?(event: string, handler: (...args: unknown[]) => void): void;
}

interface MpegtsFactory {
  Events?: {
    ERROR?: string;
  };
  isSupported(): boolean;
  createPlayer(mediaDataSource: { type: 'flv'; url: string; isLive: boolean }, config?: Record<string, unknown>): MpegtsPlayer;
}

interface UpstreamStatus {
  name: string;
  url: string;
  healthUrl: string;
  connected: boolean;
  error?: string;
}

interface ServiceSummary {
  category: string;
  total: number;
  status: Record<string, number>;
}

interface GatewayHealth {
  ok: boolean;
  inputUrl: string;
  streamUrl: string;
  bufferBytes: number;
  atccService: string;
  ptzService: string;
  upstreams: UpstreamStatus[];
  services: ServiceSummary[];
}

interface ServicesResponse {
  upstreams: UpstreamStatus[];
  services: ServiceSummary[];
}

interface SystemHealthRow {
  name: string;
  mode: string;
  status: 'online' | 'warning' | 'offline';
  url: string;
  total: number;
  connected: number;
  warning: number;
  disconnected: number;
  message: string;
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
export class DashboardComponent implements AfterViewInit, OnDestroy, OnInit {
  @ViewChildren('flvVideo') private flvVideos?: QueryList<ElementRef<HTMLVideoElement>>;

  private readonly dashboardDataService = inject(DashboardDataService);
  private readonly dashboardData = this.dashboardDataService.getDashboardData();
  private readonly flvPlayers = new Map<HTMLVideoElement, MpegtsPlayer>();
  private readonly flvStartTimers = new Map<HTMLVideoElement, number>();
  private readonly flvErrorHandlers = new Map<HTMLVideoElement, { event: string; handler: (...args: unknown[]) => void }>();
  private readonly bufferedPlaybackStarted = new WeakSet<HTMLVideoElement>();
  private readonly streamWarnings = new Set<string>();
  private readonly streamOnline = new Set<string>();
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
  systemHealthUpdatedAt = '';
  systemHealthLoading = false;
  systemHealthError = '';
  gatewayHealth: GatewayHealth | null = null;
  systemHealthRows: SystemHealthRow[] = [];
  private systemHealthTimer?: number;

  get systemHealthTotals(): { total: number; online: number; warning: number; offline: number } {
    return this.systemHealthRows.reduce(
      (totals, row) => {
        totals.total += 1;
        totals[row.status] += 1;
        return totals;
      },
      { total: 0, online: 0, warning: 0, offline: 0 },
    );
  }

  ngOnInit(): void {
    void this.refreshSystemHealth();
    this.systemHealthTimer = window.setInterval(() => {
      if (this.activeBody === 'system-health') {
        void this.refreshSystemHealth();
      }
    }, 10000);
  }

  ngAfterViewInit(): void {
    this.setupFlvPlayers();
    this.flvVideos?.changes.subscribe(() => this.setupFlvPlayers());
  }

  ngOnDestroy(): void {
    this.destroyFlvPlayers();
    if (this.systemHealthTimer !== undefined) {
      window.clearInterval(this.systemHealthTimer);
    }
  }

  selectMenu(item: MenuItem): void {
    if (!item.children?.length) {
      if (item.route !== undefined) {
        this.activeBody = item.route;
        if (item.route === 'system-health') {
          void this.refreshSystemHealth();
        }
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

  isStreamWarning(feed: { label: string; streamType: string; streamUrl?: string }): boolean {
    return feed.streamType === 'placeholder' || !feed.streamUrl || !this.streamOnline.has(feed.label) || this.streamWarnings.has(feed.label);
  }

  markStreamOnline(feedLabel: string): void {
    this.streamOnline.add(feedLabel);
    this.streamWarnings.delete(feedLabel);
  }

  markStreamWarning(feedLabel: string): void {
    this.streamOnline.delete(feedLabel);
    this.streamWarnings.add(feedLabel);
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

  async refreshSystemHealth(): Promise<void> {
    this.systemHealthLoading = true;
    this.systemHealthError = '';

    try {
      const [health, services] = await Promise.all([
        this.fetchGatewayJSON<GatewayHealth>('/healthz'),
        this.fetchGatewayJSON<ServicesResponse>('/api/services'),
      ]);
      this.gatewayHealth = health;
      this.systemHealthRows = this.buildSystemHealthRows(services, health);
      this.systemHealthUpdatedAt = new Intl.DateTimeFormat('en-IN', {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
      }).format(new Date());
    } catch (error) {
      this.systemHealthError = error instanceof Error ? error.message : 'Unable to load system health';
    } finally {
      this.systemHealthLoading = false;
    }
  }

  private buildSystemHealthRows(response: ServicesResponse, health: GatewayHealth): SystemHealthRow[] {
    const upstreams = response.upstreams ?? health.upstreams ?? [];
    const services = response.services ?? health.services ?? [];
    const upstreamByName = new Map(upstreams.map(upstream => [this.normalizeServiceName(upstream.name), upstream]));
    const rows: SystemHealthRow[] = [
      {
        name: 'Gateway',
        mode: 'Server',
        status: 'online',
        url: window.location.port === '8080' ? window.location.origin : 'http://localhost:8080',
        total: 1,
        connected: 1,
        warning: health.ok ? 0 : 1,
        disconnected: 0,
        message: health.ok ? 'Online and accepting client registrations' : 'Online, some clients need attention',
      },
    ];

    for (const summary of services) {
      const name = this.displayServiceName(summary.category);
      const upstream = upstreamByName.get(this.normalizeServiceName(name));
      const connected = Number(summary.status?.['connected'] ?? 0);
      const warning = Number(summary.status?.['warning'] ?? 0);
      const disconnected = Number(summary.status?.['disconnected'] ?? 0);
      const isSplitClient = !!upstream;
      const isOnline = isSplitClient ? upstream.connected : true;
      const status = isOnline ? 'online' : 'offline';
      const deviceNote = this.formatDeviceNote(connected, warning, disconnected);

      rows.push({
        name,
        mode: isSplitClient ? 'Client' : 'Internal',
        status,
        url: upstream?.url ?? 'gateway internal',
        total: Number(summary.total ?? 0),
        connected,
        warning,
        disconnected,
        message: upstream?.error
          ? `Not connected with gateway: ${upstream.error}`
          : (isSplitClient ? `Connected with gateway. ${deviceNote}` : `Running inside gateway. ${deviceNote}`),
      });
    }

    return rows;
  }

  private formatDeviceNote(connected: number, warning: number, disconnected: number): string {
    const notes = [`${connected} online`];
    if (warning > 0) {
      notes.push(`${warning} warning`);
    }
    if (disconnected > 0) {
      notes.push(`${disconnected} offline`);
    }
    return notes.join(', ');
  }

  private async fetchGatewayJSON<T>(path: string): Promise<T> {
    const localURL = path;
    const gatewayURL = `http://localhost:8080${path}`;
    const urls = window.location.port === '8080' ? [localURL] : [localURL, gatewayURL];
    let lastError: unknown;

    for (const url of urls) {
      try {
        const response = await fetch(url, { headers: { Accept: 'application/json' } });
        if (!response.ok) {
          throw new Error(`${url} returned ${response.status}`);
        }
        return await response.json() as T;
      } catch (error) {
        lastError = error;
      }
    }

    throw lastError instanceof Error ? lastError : new Error('Gateway health API unavailable');
  }

  private displayServiceName(category: string): string {
    const labels: Record<string, string> = {
      atcc: 'ATCC',
      vids: 'VIDS',
      'ptz-cameras': 'PTZ Camera',
      'cctv-cameras': 'CCTV Camera',
      met: 'MET',
      vms: 'VMS',
      vsds: 'VSDS',
    };
    return labels[category] ?? category.toUpperCase();
  }

  private normalizeServiceName(name: string): string {
    return name.toLowerCase().replace(/[^a-z0-9]/g, '');
  }

  applyLiveBuffer(event: Event, bufferSeconds = 2): void {
    const video = event.target as HTMLVideoElement;
    const feedLabel = video.dataset['feedLabel'];

    if (feedLabel) {
      this.markStreamOnline(feedLabel);
    }

    const liveEdge = this.getLiveEdge(video);

    if (liveEdge !== null) {
      video.currentTime = Math.max(0, liveEdge - bufferSeconds);
    }

    this.playWhenBuffered(event, bufferSeconds);
  }

  playWhenBuffered(event: Event, bufferSeconds = 5): void {
    const video = event.target as HTMLVideoElement;
    const feedLabel = video.dataset['feedLabel'];
    const requiredBuffer = Math.max(1, bufferSeconds);

    if (this.bufferedPlaybackStarted.has(video) || this.getBufferedAhead(video) < requiredBuffer) {
      return;
    }

    this.bufferedPlaybackStarted.add(video);
    if (feedLabel) {
      this.markStreamOnline(feedLabel);
    }

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
    const feedLabel = video.dataset['feedLabel'];
    const liveEdge = this.getLiveEdge(video);

    if (liveEdge === null) {
      return;
    }

    if (feedLabel) {
      this.markStreamOnline(feedLabel);
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
    const feedLabel = video.dataset['feedLabel'];

    if (!streamUrl) {
      return;
    }

    const mpegts = await this.loadMpegts();

    if (!mpegts?.isSupported()) {
      if (feedLabel) {
        this.markStreamWarning(feedLabel);
      }
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

    const errorEvent = mpegts.Events?.ERROR;
    if (errorEvent && player.on) {
      const handler = () => {
        if (feedLabel) {
          this.markStreamWarning(feedLabel);
        }
      };
      player.on(errorEvent, handler);
      this.flvErrorHandlers.set(video, { event: errorEvent, handler });
    }

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
    const errorHandler = this.flvErrorHandlers.get(video);

    if (timer !== undefined) {
      window.clearTimeout(timer);
      this.flvStartTimers.delete(video);
    }

    if (errorHandler && player.off) {
      player.off(errorHandler.event, errorHandler.handler);
      this.flvErrorHandlers.delete(video);
    }

    player.unload();
    player.detachMediaElement();
    player.destroy();
    this.flvPlayers.delete(video);
  }

}
