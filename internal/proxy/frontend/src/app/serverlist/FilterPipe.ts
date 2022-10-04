import { Pipe, PipeTransform } from "@angular/core";

@Pipe({
  name: "searchServer",
})
export class FilterPipe implements PipeTransform {
  transform(value: string[], input: string): string[] {
    if (input) {
      input = input.toLowerCase();
      return value.filter((item: string) => item.toLowerCase().includes(input));
    }
    return value;
  }
}
