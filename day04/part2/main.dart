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
    return keyVals.entries.every((o) {
        var result = validPair(o.key, o.value);
        // if (!result) print('bad key=${o.key}, val=${o.value}');
        return result;
    });
  }

  static bool validPair(String key, String val) {
    switch(key) {
      case 'byr':
      return validNum(val, 1920, 2002);
      case 'iyr':
      return validNum(val, 2010, 2020);
      case 'eyr':
      return validNum(val, 2020, 2030);
      case 'hgt':
      if (val.endsWith('cm') && validNum(val.substring(0, val.length-2), 150, 193)) return true;
      if (val.endsWith('in') && validNum(val.substring(0, val.length-2), 59,76)) return true;
      return false;
      case 'hcl':
      return validRgb(val);
      case 'ecl':
      return RegExp(r'amb|blu|brn|gry|grn|hzl|oth').hasMatch(val);
      case 'pid':
      return validPid(val);
    }
    return true;
  }

  static bool validNum(String val, int lo, int hi) {
    var v = int.parse(val);
    return v >= lo && v <= hi;
  }

  static bool validPid(String val) {
    return RegExp(r'^\d{9}$').hasMatch(val);
  }
  
  static bool validRgb(String val) {
    return RegExp(r'^#[0-9a-f]{6}$').hasMatch(val);
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
