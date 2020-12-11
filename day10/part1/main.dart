// -*- compile-command: "dart main.dart ../example1.txt ../example2.txt ../input.txt"; -*-

import 'dart:io';

main(List<String> args) {
  for (var arg in args) {
    process(arg);
  }

  print('Done.');
}

process(String filename) {
  print('Processing $filename ...');

  var ints = List<int>();
  var lines = File(filename).readAsLinesSync();
  for (var line in lines) {
    var n = int.parse(line);
    ints.add(n);
  }

  ints.sort();

  var oneDiffs = 1;
  var threeDiffs = 1;

  for (var i = 1; i < ints.length; i++) {
    var diff = ints[i] - ints[i-1];
    if (diff == 1) oneDiffs++;
    else if (diff == 3) threeDiffs++;
  }

  print('oneDiffs=$oneDiffs, threeDiffs=$threeDiffs, Solution: ${oneDiffs*threeDiffs}');
}
