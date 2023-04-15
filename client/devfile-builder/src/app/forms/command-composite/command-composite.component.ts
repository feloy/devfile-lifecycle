import { Component, EventEmitter, Output } from '@angular/core';

@Component({
  selector: 'app-command-composite',
  templateUrl: './command-composite.component.html',
  styleUrls: ['./command-composite.component.css']
})
export class CommandCompositeComponent {
  @Output() canceled = new EventEmitter<void>();

  cancel() {
    this.canceled.emit();
  }
}
