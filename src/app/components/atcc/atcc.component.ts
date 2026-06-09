import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { Component, inject, Input, OnChanges, OnInit, SimpleChanges } from '@angular/core';
import { AtccDataService, AtccSection } from '../../services/atcc-data.service';

@Component({
  selector: 'app-atcc',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './atcc.component.html',
  styleUrls: ['./atcc.component.scss'],
})
export class ATCCComponent implements OnInit, OnChanges {
  @Input() embedded = false;
  @Input() section: AtccSection | null = null;

  private readonly route = inject(ActivatedRoute);
  private readonly atccDataService = inject(AtccDataService);
  private readonly atccData = this.atccDataService.getAtccData();

  readonly updatedAt = this.atccData.updatedAt;
  readonly summaryCards = this.atccData.summaryCards;
  readonly vehicleClasses = this.atccData.vehicleClasses;
  readonly directions = this.atccData.directions;
  readonly records = this.atccData.records;
  readonly reports = this.atccData.reports;
  readonly totalVehicles = this.vehicleClasses.reduce((sum, item) => sum + item.count, 0);
  readonly pieSlices = this.buildPieSlices(this.vehicleClasses);
  activeSection: AtccSection = 'visualization';
  selectedReportIndex = 1;

  readonly subMenu: Array<{ id: AtccSection; label: string; icon: string }> = [
    { id: 'visualization', label: 'Visualization', icon: 'chart' },
    { id: 'table', label: 'Table', icon: 'grid' },
    { id: 'reports', label: 'Reports', icon: 'file' },
  ];

  readonly sectionSubMenu: Record<AtccSection, Array<{ id: string; label: string }>> = {
    visualization: [
      { id: 'overview', label: 'Overview' },
      { id: 'trends', label: 'Trends' },
      { id: 'cameras', label: 'Cameras' },
    ],
    table: [
      { id: 'all-records', label: 'All Records' },
      { id: 'speeds', label: 'Speed Data' },
      { id: 'alerts', label: 'Alerts' },
    ],
    reports: [
      { id: 'daily', label: 'Daily' },
      { id: 'weekly', label: 'Weekly' },
      { id: 'monthly', label: 'Monthly' },
    ],
  };

  activeSubSection = this.sectionSubMenu[this.activeSection][0].id;

  ngOnInit(): void {
    if (this.section && this.subMenu.some(item => item.id === this.section)) {
      this.setSection(this.section);
    }

    if (!this.embedded) {
      this.route.paramMap.subscribe(params => {
        const section = params.get('section') as AtccSection | null;
        if (section && this.subMenu.some(item => item.id === section)) {
          this.setSection(section);
          return;
        }

        this.setSection('visualization');
      });
    }
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['section'] && this.embedded && changes['section'].currentValue) {
      const section = changes['section'].currentValue as AtccSection;
      if (this.subMenu.some(item => item.id === section)) {
        this.setSection(section);
      }
    }
  }

  setSection(section: AtccSection): void {
    this.activeSection = section;
    this.activeSubSection = this.sectionSubMenu[section][0].id;
  }

  setSubSection(subSection: string): void {
    this.activeSubSection = subSection;
  }

  get selectedReport() {
    return this.reports[this.selectedReportIndex] ?? this.reports[0];
  }

  setReport(index: string): void {
    const nextIndex = Number(index);
    if (Number.isInteger(nextIndex) && this.reports[nextIndex]) {
      this.selectedReportIndex = nextIndex;
    }
  }

  getDirectionTotal(item: { east: number; north: number; south: number; west: number }): number {
    return item.east + item.north + item.south + item.west;
  }

  private buildPieSlices(items: Array<{ name: string; count: number; percent: number; color: string }>) {
    let startAngle = 0;
    return items.map(item => {
      const midAngle = startAngle + item.percent * 1.8;
      const angleRad = (midAngle - 90) * Math.PI / 180;
      const radius = 0.62;
      const labelLeft = 50 + Math.cos(angleRad) * radius * 50;
      const labelTop = 50 + Math.sin(angleRad) * radius * 50;
      const slice = {
        name: item.name,
        percent: item.percent,
        color: item.color,
        bg: this.hexToRgba(item.color, 0.82),
        left: `${labelLeft}%`,
        top: `${labelTop}%`,
      };
      startAngle += item.percent * 3.6;
      return slice;
    });
  }

  private hexToRgba(hex: string, alpha = 1) {
    const sanitized = hex.replace('#', '');
    const bigint = parseInt(sanitized, 16);
    const r = (bigint >> 16) & 255;
    const g = (bigint >> 8) & 255;
    const b = bigint & 255;
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
  }
}
