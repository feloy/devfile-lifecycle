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
import { MatSelectModule } from '@angular/material/select';
import { MatTabsModule } from '@angular/material/tabs';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';

import { AppComponent } from './app.component';
import { MetadataComponent } from './forms/metadata/metadata.component';
import { DevEnvComponent } from './forms/dev-env/dev-env.component';
import { NewDevEnvComponent } from './forms/new-dev-env/new-dev-env.component';
import { MultiTextComponent } from './controls/multi-text/multi-text.component';
import { NewUserCommandComponent } from './forms/new-user-command/new-user-command.component';
import { ContainersComponent } from './tabs/containers/containers.component';
import { ContainerComponent } from './forms/container/container.component';
import { CommandsComponent } from './tabs/commands/commands.component';
import { CommandExecComponent } from './forms/command-exec/command-exec.component';
import { CommandApplyComponent } from './forms/command-apply/command-apply.component';
import { CommandCompositeComponent } from './forms/command-composite/command-composite.component';
import { SelectContainerComponent } from './controls/select-container/select-container.component';
import { ResourcesComponent } from './tabs/resources/resources.component';
import { ResourceComponent } from './forms/resource/resource.component';
import { ImagesComponent } from './tabs/images/images.component';
import { ImageComponent } from './forms/image/image.component';
import { CommandImageComponent } from './forms/command-image/command-image.component';

@NgModule({
  declarations: [
    AppComponent,
    MetadataComponent,
    DevEnvComponent,
    NewDevEnvComponent,
    MultiTextComponent,
    NewUserCommandComponent,
    ContainersComponent,
    ContainerComponent,
    CommandsComponent,
    CommandExecComponent,
    CommandApplyComponent,
    CommandCompositeComponent,
    SelectContainerComponent,
    ResourcesComponent,
    ResourceComponent,
    ImagesComponent,
    ImageComponent,
    CommandImageComponent,
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
    MatSelectModule,
    MatTabsModule,
    MatToolbarModule,
    MatTooltipModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
