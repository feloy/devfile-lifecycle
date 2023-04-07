import { Component, Input } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';

@Component({
  selector: 'app-select-container',
  templateUrl: './select-container.component.html',
  styleUrls: ['./select-container.component.css'],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      multi: true,
      useExisting: SelectContainerComponent
    }
  ]
})
export class SelectContainerComponent implements ControlValueAccessor {
  
  @Input() containers: string[] = [];

  container: string = "";

  onChange = (_: string) => {};

  writeValue(value: any) {
    this.container = value;
  }

  registerOnChange(onChange: any) {
    this.onChange = onChange;
  }

  registerOnTouched(_: any) {}

  onSelectChange(v: string) {
    this.onChange(v);
  }
}
