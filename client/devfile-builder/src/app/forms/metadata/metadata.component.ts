import { Component } from '@angular/core';
import { FormControl } from '@angular/forms';
import { StateService } from 'src/app/state.service';
import { WasmGoService } from 'src/app/wasm-go.service';

@Component({
  selector: 'app-metadata',
  templateUrl: './metadata.component.html',
  styleUrls: ['./metadata.component.css']
})
export class MetadataComponent {
  name = new FormControl('');
  displayName = new FormControl('');
  description = new FormControl('');

  constructor(
    private wasm: WasmGoService,
    private state: StateService,
  ) {}

  onSave() {
    const newDevfile = this.wasm.setMetadata({
      name: this.name.value,
      displayName: this.displayName.value,
      description: this.description.value,
    })
    this.state.changeDevfileYaml(newDevfile);
  }
}
