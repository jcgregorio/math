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

  return mb;
}();
