import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ReactiveFormsModule } from '@angular/forms';

import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatTabsModule } from '@angular/material/tabs';
import { MatToolbarModule } from '@angular/material/toolbar';

import { AppComponent } from './app.component';
import { MetadataComponent } from './forms/metadata/metadata.component';
import { DevEnvComponent } from './forms/dev-env/dev-env.component';
import { NewDevEnvComponent } from './forms/new-dev-env/new-dev-env.component';
import { MultiTextComponent } from './controls/multi-text/multi-text.component';
import { NewUserCommandComponent } from './forms/new-user-command/new-user-command.component';

@NgModule({
  declarations: [
    AppComponent,
    MetadataComponent,
    DevEnvComponent,
    NewDevEnvComponent,
    MultiTextComponent,
    NewUserCommandComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    ReactiveFormsModule,

    MatButtonModule,
    MatCardModule,
    MatCheckboxModule,
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
