import { Component } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { StateService } from 'src/app/services/state.service';
import { WasmGoService } from 'src/app/services/wasm-go.service';

@Component({
  selector: 'app-image',
  templateUrl: './image.component.html',
  styleUrls: ['./image.component.css']
})
export class ImageComponent {
  form: FormGroup;

  constructor(
    private wasm: WasmGoService,
    private state: StateService,
  ) {
    this.form = new FormGroup({
      name: new FormControl(""),
      imageName: new FormControl(""),
      args: new FormControl([]),
      buildContext: new FormControl(""),
      rootRequired: new FormControl(false),
      uri: new FormControl(""),
    })
  }

  create() {
    const newDevfile = this.wasm.addImage(this.form.value);
    this.state.changeDevfileYaml(newDevfile);
  }

}
