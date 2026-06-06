import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTableModule } from '@angular/material/table';

import { ATCCComponent } from './atcc.component';

@NgModule({
  declarations: [ATCCComponent],
  imports: [
    CommonModule,
    MatTableModule   // ✅ Add this
  ]
})
export class ATCCModule {}
