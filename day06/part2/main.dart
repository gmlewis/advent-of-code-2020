// -*- compile-command: "dart main.dart ../input.txt"; -*-

import 'dart:io';

main(List<String> args) {
  for (var arg in args) {
    process(arg);
  }

  print('Done.');
}

process(String filename) {
  print('Processing $filename ...');

  var buf = File(filename).readAsStringSync();
  var groups = buf.split('\n\n');

  var count = 0;
  for (var group in groups) {
    group = group.trim();
    var runes = Map<int, int>();
    var lines = group.split('\n');
    lines.forEach((line) {
        line.runes.forEach((r) {
            runes.update(r, (v)=>v+1, ifAbsent: ()=>1);
            if (runes[r] == lines.length) count++;
        });
    });
  }
  print('Solution: $count');
}
