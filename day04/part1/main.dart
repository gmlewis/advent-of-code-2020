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

  var passports = Passport.read(filename);

  print('${passports.length} valid passports');
}

class Passport {
  Map<String, String> keyVals;

  static const required = ['byr', 'iyr', 'eyr', 'hgt', 'hcl', 'ecl', 'pid'];
  
  Passport(String buf) {
    buf = buf.replaceAll('\n', ' ');
    this.keyVals = {};

    var fields = buf.split(' ');
    for (var keyVal in fields) {
      if (keyVal.isEmpty) continue;
      var parts = keyVal.split(':');
      this.keyVals[parts[0]] = parts[1];
    }
  }

  bool valid() {
    for (var req in required) {
      if (!keyVals.containsKey(req)) return false;
    }
    return true;
  }
  
  static List<Passport> read(String filename) {
    var passports = List<Passport>();
    var buf = File(filename).readAsStringSync();
    var groups = buf.split('\n\n');
    for (var group in groups) {
      var passport = Passport(group);
      if (passport.valid()) {
        passports.add(passport);
      }
    }
    return passports;
  }
}
