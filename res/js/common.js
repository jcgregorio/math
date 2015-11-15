/**
 * $$ returns a real JS array of DOM elements that match the CSS query selector.
 *
 * A shortcut for jQuery-like $ behavior.
 **/
function $$(query, ele) {
  if (!ele) {
    ele = document;
  }
  return Array.prototype.map.call(ele.querySelectorAll(query), function(e) { return e; });
}


/**
 * $$$ returns the DOM element that match the CSS query selector.
 *
 * A shortcut for document.querySelector.
 **/
function $$$(query, ele) {
  if (!ele) {
    ele = document;
  }
  return ele.querySelector(query);
}


window.mb = window.mb || function() {
  var mb = {};

  /**
   * clearChildren removes all children of the passed in node.
   */
  mb.clearChildren = function(ele) {
    while (ele.firstChild) {
      ele.removeChild(ele.firstChild);
    }
  }


  mb.rint = function(min, max) {
    return Math.floor(Math.random() * (max - min)) + min;
  };

  mb.plusminus = function() {
    if (Math.random() > 0.5) {
      return 1;
    } else {
      return -1;
    }
  };

  mb.disableAll = function(children) {
    for (var i = children.length - 1; i >= 0; i--) {
      var c = children[i];
      if (c.hasChildNodes()) {
        mb.disableAll(c.children);
      }
      if ('disabled' in c) {
        c.disabled = true;
      }
    }
  };

  // Returns a Promise that uses XMLHttpRequest to make a request with the given
  // method to the given URL with the given headers and body.
  mb.request = function(method, url, body, headers) {
    // Return a new promise.
    return new Promise(function(resolve, reject) {
      // Do the usual XHR stuff
      var req = new XMLHttpRequest();
      req.open(method, url);
      if (headers) {
        for (var k in headers) {
          req.setRequestHeader(k, headers[k]);
        }
      }

      req.onload = function() {
        // This is called even on 404 etc
        // so check the status
        if (req.status == 200) {
          // Resolve the promise with the response text
          resolve(req.response);
        } else {
          // Otherwise reject with the status text
          // which will hopefully be a meaningful error
          reject(req.response);
        }
      };

      // Handle network errors
      req.onerror = function() {
        reject(Error("Network Error"));
      };

      // Make the request
      req.send(body);
    });
  }

  // Returns a Promise that uses XMLHttpRequest to make a request to the given URL.
  mb.get = function(url) {
    return mb.request('GET', url);
  }


  // Returns a Promise that uses XMLHttpRequest to make a POST request to the
  // given URL with the given JSON body.
  mb.post = function(url, body) {
    return mb.request('POST', url, body, {"Content-Type": "application/json"});
  }

  // Returns a Promise that uses XMLHttpRequest to make a DELETE request to the
  // given URL.
  mb.delete = function(url) {
    return mb.request('DELETE', url);
  }

  return mb;
}();
