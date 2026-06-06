import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { DashboardComponent } from './components/dashboard/dashboard.component';
import { ATCCComponent } from './components/atcc/atcc.component';
import { VIDSComponent } from './components/vids/vids.component';
import { SpeedViolationComponent } from './components/speed-violation/speed-violation.component';
import { DeviceStatusMapComponent } from './components/device-status-map/device-status-map.component';
import { DeviceStatusOverviewComponent } from './components/device-status-overview/device-status-overview.component';


const routes: Routes = [
  { path: '', component: DashboardComponent },          // default route
  { path: 'atcc', component: ATCCComponent },
  { path: 'vids', component: VIDSComponent },
  { path: 'speed-violation', component: SpeedViolationComponent },
  { path: 'device-map', component: DeviceStatusMapComponent },
  { path: 'device-overview', component: DeviceStatusOverviewComponent },
  { path: '**', redirectTo: '' }                        // wildcard fallback
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
