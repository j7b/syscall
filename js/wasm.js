'use strict';

const WASM_URL = 'main.wasm';

var wasm;
var logLine = [];
var memory8;

var rc = 8;
var refs = new Map();
refs.set(8,window);

var decode = new TextDecoder("utf-8").decode;

function dattype(o) {
  if (typeof o == "undefined") {
    return 0;
  }
  if (o == null) {
    return 1;
  }
  var t = typeof o;
  switch (t) {
    case "object":
    return 6;
    case "boolean":
    return 2;
    case "number":
    return 3;
    case "string":
    return 4;
    case "symbol":
    return 5;
    case "function":
    return 7;
  }
}

var importObject = {
  env: {
    'main.newRef':function() {
      rc++;
      return rc;
    },
    'main.readRefString':function(r,offset) {
      var o = refs.get(r);
      if (typeof o.charCodeAt === undefined) {
        o = o.toString();
      }
      var s = o.charCodeAt(offset);
      if (isNaN(s)) {
        return -1;
      }
      return s;
    },
    'main.typeOf':dattype,
    'main.getRef':function(r,vr) {
      var thing = refs.get(r);
      var key = refs.get(vr);
      var obj = Reflect.get(thing,key);
      if (typeof obj == 'undefined') {
        return -2;
      }
      if (obj == null) {
        return -1;
      }
      rc++;
      refs.set(rc,obj);
      return rc;
    },
    'main.setRef':function(r,v) {
      refs.set(r,v);
    },
    'main.appendRefString':function(r,c) {
      var s = refs.get(r);
      if (!s) {
        s = "";
      }
      s += String.fromCharCode(c);
      refs.set(r,s);
        },
    'github.com/j7b/syscall/js.refRelease':function(r) {
      refs.delete(r);
      if (r == rc) {
        rc--;
      }
    },
    'main.intRef':function(r) {
      return refs.get(r);
    },
    'github.com/j7b/syscall/js.newBool':function(r) {
      rc++;
      refs.set(rc,r==1);
      return rc;
    },
    'github.com/j7b/syscall/js.newInt':function(r) {
      rc++;
      refs.set(rc,r);
      return rc;
    },
    'github.com/j7b/syscall/js.newFloat':function(r) {
      rc++;
      refs.set(rc,r);
      return rc;
    },
    'github.com/j7b/syscall/js.newString':function() {
      rc++;
      refs.set(rc,"");
      return rc;
    },
    'github.com/j7b/syscall/js.newArray':function() {
      rc++;
      refs.set(rc,[]);
      return rc;
    },
    'github.com/j7b/syscall/js.appendString':function(ref,rune) {
      var s = refs.get(ref);
      s += String.fromCharCode(rune);
      refs.set(ref,s);
    },
    'github.com/j7b/syscall/js.intAt':function(ref,index) {
      var s = refs.get(ref);
      return s[index];
    },
    'github.com/j7b/syscall/js.typeOf':function(ref) {
      var s = refs.get(ref);
      return dattype(s);
    },
    'github.com/j7b/syscall/js.getInt':function(ref) {
      var s = refs.get(ref);
      return s;
    },
    'github.com/j7b/syscall/js.setInt':function(ref,v) {
      refs.set(ref,v);
    },
    'github.com/j7b/syscall/js.getFloat':function(ref) {
      var s = refs.get(ref);
      return s;
    },
    'github.com/j7b/syscall/js.setFloat':function(ref,v) {
      refs.set(ref,v);
    },
    'github.com/j7b/syscall/js.setVal':function(ref,key,val) {
      ref = refs.get(ref);
      key = refs.get(key);
      val = refs.get(val);
      Reflect.set(ref,key,val);
    }, 
    'github.com/j7b/syscall/js.getRef':function(ref,k) {
      var key = refs.get(k);
      ref = refs.get(ref);
      var val = Reflect.get(ref,key);
      rc++;
      refs.set(rc,val);
      return rc;
    },
    'github.com/j7b/syscall/js.getIndex':function(ref,index) {
      var arr = refs.get(ref);
      var val = arr[index];
      rc++;
      refs.set(rc,val);
      return rc;
    },
    'github.com/j7b/syscall/js.charCodeAt':function(ref,index) {
      ref = refs.get(ref);
      return ref.charCodeAt(index);
    },
    'github.com/j7b/syscall/js.setIndex':function(ref,index,val) {
      var arr = refs.get(ref);
      arr[index] = refs.get(val);
      refs.set(ref,arr);
    },
    'github.com/j7b/syscall/js.pushRef':function(ref,val) {
      var arr = refs.get(ref);
      arr.push(refs.get(val));
      refs.set(ref,arr);
    },
    'github.com/j7b/syscall/js.instanceOf':function(r1,r2) {
      r1 = refs.get(r1);
      r2 = refs.get(r2);
      if (r1 instanceof r2) {
        return 1;
      }
      return 0;
    },
    'github.com/j7b/syscall/js.lengthOf':function(ref) {
      ref = refs.get(ref);
      return ref.length;
    },
    'github.com/j7b/syscall/js.call':function(ref,key,arg) {
      ref = refs.get(ref);
      key = refs.get(key);
      if (arg > 0) {
        arg = refs.get(arg);
      } else {
        arg = undefined;
      }
      var f = ref[key];
      var out = Reflect.apply(f,ref,arg);
      rc++;
      refs.set(rc,out);
      return rc;
    },
    'github.com/j7b/syscall/js.invoke':function(ref,arg) {
      ref = refs.get(ref);
      if (arg > 0) {
        arg = refs.get(arg);
      } else {
        arg = undefined;
      }
      var out = Reflect.apply(ref,undefined,arg);
      rc++;
      refs.set(rc,out);
      return rc;
    },
    'github.com/j7b/syscall/js.noo':function(ref,arg) {
      ref = refs.get(ref);
      if (arg > 0) {
        arg = refs.get(arg);
      } else {
        arg = undefined;
      }
      var out = Reflect.construct(ref,arg);
      rc++;
      refs.set(rc,out);
      return rc;
    },
    'github.com/j7b/syscall/js.zero':function(ref,t) {
      var val;
      switch (t) {
        case 0:
        refs.set(ref,undefined);
        return;
        case 1:
        refs.set(ref,null);
        return;
        case 2:
        refs.set(ref,false);
        return;
        case 3:
        refs.set(ref,0);
        return;
        case 4:
        refs.set(ref,"");
        return;
        case 5:
        refs.set(ref,Symbol());
        return;
        case 6:
        refs.set(ref,{});
        return;
        case 7:
        refs.set(ref,function(){});
        return;
      }
    },
    io_get_stdout: function() {
      return 1;
    },
    resource_write: function(fd, ptr, len) {
      if (fd == 1) {
        for (let i=0; i<len; i++) {
          let c = memory8[ptr+i];
          if (c == 13) { // CR
            // ignore
          } else if (c == 10) { // LF
            // write line
            let line = new TextDecoder("utf-8").decode(new Uint8Array(logLine));
            logLine = [];
            console.log(line);
          } else {
            logLine.push(c);
          }
        }
      } else {
        console.error('invalid file descriptor:', fd);
      }
    },
  },
};

function init() {
  const derp = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
  derp(fetch(WASM_URL), importObject).then(function(obj) {
    wasm = obj.instance;
    memory8 = new Uint8Array(wasm.exports.memory.buffer);
    wasm.exports.cwa_main();
  })
}

init();
