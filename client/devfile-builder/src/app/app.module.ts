import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ReactiveFormsModule } from '@angular/forms';

import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatTabsModule } from '@angular/material/tabs';
import { MatToolbarModule } from '@angular/material/toolbar';


import { AppComponent } from './app.component';
import { MetadataComponent } from './forms/metadata/metadata.component';

@NgModule({
  declarations: [
    AppComponent,
    MetadataComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    ReactiveFormsModule,

    MatButtonModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatTabsModule,
    MatToolbarModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
