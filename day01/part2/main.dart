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
    if (find(v, vals)) {
      break;
    }
    vals.add(sum-v);
  }
}

bool find(int v, Set<int> vals) {
	for (var d2 in vals) {
		var v2 = sum - d2;
		if (vals.contains(v+v2)) {
      print('$v + $v2 + ${sum-v-v2} = $sum\n$v * $v2 * ${sum-v-v2} = ${v * v2 * (sum-v-v2)}');
			return true;
		}
	}
	return false;
}
