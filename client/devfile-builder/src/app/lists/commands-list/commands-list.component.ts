import { Component, Input } from '@angular/core';
import { StateService } from 'src/app/services/state.service';
import { Command, WasmGoService } from 'src/app/services/wasm-go.service';

@Component({
  selector: 'app-commands-list',
  templateUrl: './commands-list.component.html',
  styleUrls: ['./commands-list.component.css']
})
export class CommandsListComponent {
  @Input() commands: Command[] | undefined;
  @Input() kind: string = "";

  constructor(
    private wasm: WasmGoService,
    private state: StateService,
  ) {}

  setDefault(command: string, group: string) {
    const newDevfile = this.wasm.setDefaultCommand(command, group);
    this.state.changeDevfileYaml(newDevfile);
  }
}
