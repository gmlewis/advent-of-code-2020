// -*- compile-command: "dart main.dart ../example1.txt ../input.txt"; -*-

import 'dart:io';

main(List<String> args) {
  for (var arg in args) {
    process(arg);
  }

  print('Done.');
}

process(String filename) {
  print('Processing $filename ...');

	var lineRE  = RegExp(r'^(\d+)-(\d+) ([a-z]): (.*)$');

  var count = 0;
  var lines = File(filename).readAsLinesSync();
  for (var line in lines) {
    var m = lineRE.allMatches(line).first;

    var start = int.parse(m.group(1));
    var end = int.parse(m.group(2));
    var letter = m.group(3);
    var passwd = m.group(4);

    if (valid(start, end, letter, passwd)) {
      count++;
    }
  }

  print('$count valid passwords');
}

bool valid(int start, int end, String letter, String passwd) {
  // print('$start-$end $letter: $passwd');
  var first = passwd.substring(start-1,start) == letter;
  var second = passwd.substring(end-1,end) == letter;
  return first != second;
}
