import { Component } from '@angular/core';
import { StateService } from 'src/app/services/state.service';
import { Command } from 'src/app/services/wasm-go.service';

@Component({
  selector: 'app-commands',
  templateUrl: './commands.component.html',
  styleUrls: ['./commands.component.css']
})
export class CommandsComponent {

  forceDisplayExecForm: boolean = false;
  forceDisplayApplyForm: boolean = false;
  forceDisplayImageForm: boolean = false;
  forceDisplayCompositeForm: boolean = false;

  commands: Command[] | undefined = [];

  constructor(
    private state: StateService,
  ) {}

  ngOnInit() {
    const that = this;
    this.state.state.subscribe(async newContent => {
      that.commands = newContent?.commands;
      if (this.commands == null) {
        return
      }
      that.forceDisplayExecForm = false;
      this.forceDisplayApplyForm = false;
      this.forceDisplayCompositeForm = false;
    });
  }

  displayExecForm() {
    this.forceDisplayExecForm = true;
  }

  displayApplyForm() {
    this.forceDisplayApplyForm = true;
  }

  displayImageForm() {
    this.forceDisplayImageForm = true;
  }

  displayCompositeForm() {
    this.forceDisplayCompositeForm = true;
  }

}
