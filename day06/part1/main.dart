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
    group = group.replaceAll('\n', '');
    var chars = Set<int>();
    group.runes.forEach((ch) => chars.add(ch));
    count += chars.length;
  }
  print('Solution: $count');
}
