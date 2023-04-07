import { Component } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { StateService } from 'src/app/services/state.service';
import { WasmGoService } from 'src/app/services/wasm-go.service';

@Component({
  selector: 'app-container',
  templateUrl: './container.component.html',
  styleUrls: ['./container.component.css']
})
export class ContainerComponent {
  form: FormGroup;

  constructor(
    private wasm: WasmGoService,
    private state: StateService,
  ) {
    this.form = new FormGroup({
      name: new FormControl(""),
      image: new FormControl(""),
      command: new FormControl([]),
      args: new FormControl([]),
    })
  }

  create() {
    const newDevfile = this.wasm.addContainer(this.form.value);
    this.state.changeDevfileYaml(newDevfile);
  }
}
