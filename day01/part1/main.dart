// -*- compile-command: "dart main.dart ../input.txt"; -*-

import 'dart:io';

const sum = 2020;

main(List<String> args) {
  for (var arg in args) {
    process(arg);
  }

  print('Done.');
}

process(String filename) {
  print('Processing $filename ...');

  var vals = Set<int>();
  var lines = File(filename).readAsLinesSync();
  for (var line in lines) {
    var v = int.parse(line);
    if (vals.contains(v)) {
      print('$v + ${sum-v} = sum\n$v * ${sum-v} = ${v * (sum-v)}');
      break;
    }
    vals.add(sum-v);
  }
}
