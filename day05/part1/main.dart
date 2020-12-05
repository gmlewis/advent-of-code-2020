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

  var ids = List<int>();
  var lines = File(filename).readAsLinesSync();
  for (var line in lines) {
    var id = spaceID(line);
    ids.add(id);
  }

  ids.sort();

  print('Solution: ${ids[ids.length-1]}');
}

int spaceID(String s) {
  s = s.replaceAll('F', '0');
  s = s.replaceAll('B', '1');
  s = s.replaceAll('L', '0');
  s = s.replaceAll('R', '1');
  var val = int.parse(s, radix: 2);
  return val;
}
