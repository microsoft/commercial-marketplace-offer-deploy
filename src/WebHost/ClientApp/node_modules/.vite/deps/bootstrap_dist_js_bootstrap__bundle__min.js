import {
  __commonJS
} from "./chunk-UXIASGQL.js";

// node_modules/bootstrap/dist/js/bootstrap.bundle.min.js
var require_bootstrap_bundle_min = __commonJS({
  "node_modules/bootstrap/dist/js/bootstrap.bundle.min.js"(exports, module) {
    !function(t, e) {
      "object" == typeof exports && "undefined" != typeof module ? module.exports = e() : "function" == typeof define && define.amd ? define(e) : (t = "undefined" != typeof globalThis ? globalThis : t || self).bootstrap = e();
    }(exports, function() {
      "use strict";
      const t = /* @__PURE__ */ new Map(), e = { set(e2, i2, n2) {
        t.has(e2) || t.set(e2, /* @__PURE__ */ new Map());
        const s2 = t.get(e2);
        s2.has(i2) || 0 === s2.size ? s2.set(i2, n2) : console.error(`Bootstrap doesn't allow more than one instance per element. Bound instance: ${Array.from(s2.keys())[0]}.`);
      }, get: (e2, i2) => t.has(e2) && t.get(e2).get(i2) || null, remove(e2, i2) {
        if (!t.has(e2))
          return;
        const n2 = t.get(e2);
        n2.delete(i2), 0 === n2.size && t.delete(e2);
      } }, i = "transitionend", n = (t2) => (t2 && window.CSS && window.CSS.escape && (t2 = t2.replace(/#([^\s"#']+)/g, (t3, e2) => `#${CSS.escape(e2)}`)), t2), s = (t2) => {
        t2.dispatchEvent(new Event(i));
      }, o = (t2) => !(!t2 || "object" != typeof t2) && (void 0 !== t2.jquery && (t2 = t2[0]), void 0 !== t2.nodeType), r = (t2) => o(t2) ? t2.jquery ? t2[0] : t2 : "string" == typeof t2 && t2.length > 0 ? document.querySelector(n(t2)) : null, a = (t2) => {
        if (!o(t2) || 0 === t2.getClientRects().length)
          return false;
        const e2 = "visible" === getComputedStyle(t2).getPropertyValue("visibility"), i2 = t2.closest("details:not([open])");
        if (!i2)
          return e2;
        if (i2 !== t2) {
          const e3 = t2.closest("summary");
          if (e3 && e3.parentNode !== i2)
            return false;
          if (null === e3)
            return false;
        }
        return e2;
      }, l = (t2) => !t2 || t2.nodeType !== Node.ELEMENT_NODE || !!t2.classList.contains("disabled") || (void 0 !== t2.disabled ? t2.disabled : t2.hasAttribute("disabled") && "false" !== t2.getAttribute("disabled")), c = (t2) => {
        if (!document.documentElement.attachShadow)
          return null;
        if ("function" == typeof t2.getRootNode) {
          const e2 = t2.getRootNode();
          return e2 instanceof ShadowRoot ? e2 : null;
        }
        return t2 instanceof ShadowRoot ? t2 : t2.parentNode ? c(t2.parentNode) : null;
      }, h = () => {
      }, d = (t2) => {
        t2.offsetHeight;
      }, u = () => window.jQuery && !document.body.hasAttribute("data-bs-no-jquery") ? window.jQuery : null, f = [], p = () => "rtl" === document.documentElement.dir, m = (t2) => {
        var e2;
        e2 = () => {
          const e3 = u();
          if (e3) {
            const i2 = t2.NAME, n2 = e3.fn[i2];
            e3.fn[i2] = t2.jQueryInterface, e3.fn[i2].Constructor = t2, e3.fn[i2].noConflict = () => (e3.fn[i2] = n2, t2.jQueryInterface);
          }
        }, "loading" === document.readyState ? (f.length || document.addEventListener("DOMContentLoaded", () => {
          for (const t3 of f)
            t3();
        }), f.push(e2)) : e2();
      }, g = (t2, e2 = [], i2 = t2) => "function" == typeof t2 ? t2(...e2) : i2, _ = (t2, e2, n2 = true) => {
        if (!n2)
          return void g(t2);
        const o2 = ((t3) => {
          if (!t3)
            return 0;
          let { transitionDuration: e3, transitionDelay: i2 } = window.getComputedStyle(t3);
          const n3 = Number.parseFloat(e3), s2 = Number.parseFloat(i2);
          return n3 || s2 ? (e3 = e3.split(",")[0], i2 = i2.split(",")[0], 1e3 * (Number.parseFloat(e3) + Number.parseFloat(i2))) : 0;
        })(e2) + 5;
        let r2 = false;
        const a2 = ({ target: n3 }) => {
          n3 === e2 && (r2 = true, e2.removeEventListener(i, a2), g(t2));
        };
        e2.addEventListener(i, a2), setTimeout(() => {
          r2 || s(e2);
        }, o2);
      }, b = (t2, e2, i2, n2) => {
        const s2 = t2.length;
        let o2 = t2.indexOf(e2);
        return -1 === o2 ? !i2 && n2 ? t2[s2 - 1] : t2[0] : (o2 += i2 ? 1 : -1, n2 && (o2 = (o2 + s2) % s2), t2[Math.max(0, Math.min(o2, s2 - 1))]);
      }, v = /[^.]*(?=\..*)\.|.*/, y = /\..*/, w = /::\d+$/, A = {};
      let E = 1;
      const T = { mouseenter: "mouseover", mouseleave: "mouseout" }, C = /* @__PURE__ */ new Set(["click", "dblclick", "mouseup", "mousedown", "contextmenu", "mousewheel", "DOMMouseScroll", "mouseover", "mouseout", "mousemove", "selectstart", "selectend", "keydown", "keypress", "keyup", "orientationchange", "touchstart", "touchmove", "touchend", "touchcancel", "pointerdown", "pointermove", "pointerup", "pointerleave", "pointercancel", "gesturestart", "gesturechange", "gestureend", "focus", "blur", "change", "reset", "select", "submit", "focusin", "focusout", "load", "unload", "beforeunload", "resize", "move", "DOMContentLoaded", "readystatechange", "error", "abort", "scroll"]);
      function O(t2, e2) {
        return e2 && `${e2}::${E++}` || t2.uidEvent || E++;
      }
      function x(t2) {
        const e2 = O(t2);
        return t2.uidEvent = e2, A[e2] = A[e2] || {}, A[e2];
      }
      function k(t2, e2, i2 = null) {
        return Object.values(t2).find((t3) => t3.callable === e2 && t3.delegationSelector === i2);
      }
      function L(t2, e2, i2) {
        const n2 = "string" == typeof e2, s2 = n2 ? i2 : e2 || i2;
        let o2 = I(t2);
        return C.has(o2) || (o2 = t2), [n2, s2, o2];
      }
      function S(t2, e2, i2, n2, s2) {
        if ("string" != typeof e2 || !t2)
          return;
        let [o2, r2, a2] = L(e2, i2, n2);
        if (e2 in T) {
          const t3 = (t4) => function(e3) {
            if (!e3.relatedTarget || e3.relatedTarget !== e3.delegateTarget && !e3.delegateTarget.contains(e3.relatedTarget))
              return t4.call(this, e3);
          };
          r2 = t3(r2);
        }
        const l2 = x(t2), c2 = l2[a2] || (l2[a2] = {}), h2 = k(c2, r2, o2 ? i2 : null);
        if (h2)
          return void (h2.oneOff = h2.oneOff && s2);
        const d2 = O(r2, e2.replace(v, "")), u2 = o2 ? function(t3, e3, i3) {
          return function n3(s3) {
            const o3 = t3.querySelectorAll(e3);
            for (let { target: r3 } = s3; r3 && r3 !== this; r3 = r3.parentNode)
              for (const a3 of o3)
                if (a3 === r3)
                  return P(s3, { delegateTarget: r3 }), n3.oneOff && N.off(t3, s3.type, e3, i3), i3.apply(r3, [s3]);
          };
        }(t2, i2, r2) : function(t3, e3) {
          return function i3(n3) {
            return P(n3, { delegateTarget: t3 }), i3.oneOff && N.off(t3, n3.type, e3), e3.apply(t3, [n3]);
          };
        }(t2, r2);
        u2.delegationSelector = o2 ? i2 : null, u2.callable = r2, u2.oneOff = s2, u2.uidEvent = d2, c2[d2] = u2, t2.addEventListener(a2, u2, o2);
      }
      function D(t2, e2, i2, n2, s2) {
        const o2 = k(e2[i2], n2, s2);
        o2 && (t2.removeEventListener(i2, o2, Boolean(s2)), delete e2[i2][o2.uidEvent]);
      }
      function $(t2, e2, i2, n2) {
        const s2 = e2[i2] || {};
        for (const [o2, r2] of Object.entries(s2))
          o2.includes(n2) && D(t2, e2, i2, r2.callable, r2.delegationSelector);
      }
      function I(t2) {
        return t2 = t2.replace(y, ""), T[t2] || t2;
      }
      const N = { on(t2, e2, i2, n2) {
        S(t2, e2, i2, n2, false);
      }, one(t2, e2, i2, n2) {
        S(t2, e2, i2, n2, true);
      }, off(t2, e2, i2, n2) {
        if ("string" != typeof e2 || !t2)
          return;
        const [s2, o2, r2] = L(e2, i2, n2), a2 = r2 !== e2, l2 = x(t2), c2 = l2[r2] || {}, h2 = e2.startsWith(".");
        if (void 0 === o2) {
          if (h2)
            for (const i3 of Object.keys(l2))
              $(t2, l2, i3, e2.slice(1));
          for (const [i3, n3] of Object.entries(c2)) {
            const s3 = i3.replace(w, "");
            a2 && !e2.includes(s3) || D(t2, l2, r2, n3.callable, n3.delegationSelector);
          }
        } else {
          if (!Object.keys(c2).length)
            return;
          D(t2, l2, r2, o2, s2 ? i2 : null);
        }
      }, trigger(t2, e2, i2) {
        if ("string" != typeof e2 || !t2)
          return null;
        const n2 = u();
        let s2 = null, o2 = true, r2 = true, a2 = false;
        e2 !== I(e2) && n2 && (s2 = n2.Event(e2, i2), n2(t2).trigger(s2), o2 = !s2.isPropagationStopped(), r2 = !s2.isImmediatePropagationStopped(), a2 = s2.isDefaultPrevented());
        const l2 = P(new Event(e2, { bubbles: o2, cancelable: true }), i2);
        return a2 && l2.preventDefault(), r2 && t2.dispatchEvent(l2), l2.defaultPrevented && s2 && s2.preventDefault(), l2;
      } };
      function P(t2, e2 = {}) {
        for (const [i2, n2] of Object.entries(e2))
          try {
            t2[i2] = n2;
          } catch (e3) {
            Object.defineProperty(t2, i2, { configurable: true, get: () => n2 });
          }
        return t2;
      }
      function M(t2) {
        if ("true" === t2)
          return true;
        if ("false" === t2)
          return false;
        if (t2 === Number(t2).toString())
          return Number(t2);
        if ("" === t2 || "null" === t2)
          return null;
        if ("string" != typeof t2)
          return t2;
        try {
          return JSON.parse(decodeURIComponent(t2));
        } catch (e2) {
          return t2;
        }
      }
      function j(t2) {
        return t2.replace(/[A-Z]/g, (t3) => `-${t3.toLowerCase()}`);
      }
      const F = { setDataAttribute(t2, e2, i2) {
        t2.setAttribute(`data-bs-${j(e2)}`, i2);
      }, removeDataAttribute(t2, e2) {
        t2.removeAttribute(`data-bs-${j(e2)}`);
      }, getDataAttributes(t2) {
        if (!t2)
          return {};
        const e2 = {}, i2 = Object.keys(t2.dataset).filter((t3) => t3.startsWith("bs") && !t3.startsWith("bsConfig"));
        for (const n2 of i2) {
          let i3 = n2.replace(/^bs/, "");
          i3 = i3.charAt(0).toLowerCase() + i3.slice(1, i3.length), e2[i3] = M(t2.dataset[n2]);
        }
        return e2;
      }, getDataAttribute: (t2, e2) => M(t2.getAttribute(`data-bs-${j(e2)}`)) };
      class H {
        static get Default() {
          return {};
        }
        static get DefaultType() {
          return {};
        }
        static get NAME() {
          throw new Error('You have to implement the static method "NAME", for each component!');
        }
        _getConfig(t2) {
          return t2 = this._mergeConfigObj(t2), t2 = this._configAfterMerge(t2), this._typeCheckConfig(t2), t2;
        }
        _configAfterMerge(t2) {
          return t2;
        }
        _mergeConfigObj(t2, e2) {
          const i2 = o(e2) ? F.getDataAttribute(e2, "config") : {};
          return { ...this.constructor.Default, ..."object" == typeof i2 ? i2 : {}, ...o(e2) ? F.getDataAttributes(e2) : {}, ..."object" == typeof t2 ? t2 : {} };
        }
        _typeCheckConfig(t2, e2 = this.constructor.DefaultType) {
          for (const [n2, s2] of Object.entries(e2)) {
            const e3 = t2[n2], r2 = o(e3) ? "element" : null == (i2 = e3) ? `${i2}` : Object.prototype.toString.call(i2).match(/\s([a-z]+)/i)[1].toLowerCase();
            if (!new RegExp(s2).test(r2))
              throw new TypeError(`${this.constructor.NAME.toUpperCase()}: Option "${n2}" provided type "${r2}" but expected type "${s2}".`);
          }
          var i2;
        }
      }
      class W extends H {
        constructor(t2, i2) {
          super(), (t2 = r(t2)) && (this._element = t2, this._config = this._getConfig(i2), e.set(this._element, this.constructor.DATA_KEY, this));
        }
        dispose() {
          e.remove(this._element, this.constructor.DATA_KEY), N.off(this._element, this.constructor.EVENT_KEY);
          for (const t2 of Object.getOwnPropertyNames(this))
            this[t2] = null;
        }
        _queueCallback(t2, e2, i2 = true) {
          _(t2, e2, i2);
        }
        _getConfig(t2) {
          return t2 = this._mergeConfigObj(t2, this._element), t2 = this._configAfterMerge(t2), this._typeCheckConfig(t2), t2;
        }
        static getInstance(t2) {
          return e.get(r(t2), this.DATA_KEY);
        }
        static getOrCreateInstance(t2, e2 = {}) {
          return this.getInstance(t2) || new this(t2, "object" == typeof e2 ? e2 : null);
        }
        static get VERSION() {
          return "5.3.2";
        }
        static get DATA_KEY() {
          return `bs.${this.NAME}`;
        }
        static get EVENT_KEY() {
          return `.${this.DATA_KEY}`;
        }
        static eventName(t2) {
          return `${t2}${this.EVENT_KEY}`;
        }
      }
      const B = (t2) => {
        let e2 = t2.getAttribute("data-bs-target");
        if (!e2 || "#" === e2) {
          let i2 = t2.getAttribute("href");
          if (!i2 || !i2.includes("#") && !i2.startsWith("."))
            return null;
          i2.includes("#") && !i2.startsWith("#") && (i2 = `#${i2.split("#")[1]}`), e2 = i2 && "#" !== i2 ? n(i2.trim()) : null;
        }
        return e2;
      }, z = { find: (t2, e2 = document.documentElement) => [].concat(...Element.prototype.querySelectorAll.call(e2, t2)), findOne: (t2, e2 = document.documentElement) => Element.prototype.querySelector.call(e2, t2), children: (t2, e2) => [].concat(...t2.children).filter((t3) => t3.matches(e2)), parents(t2, e2) {
        const i2 = [];
        let n2 = t2.parentNode.closest(e2);
        for (; n2; )
          i2.push(n2), n2 = n2.parentNode.closest(e2);
        return i2;
      }, prev(t2, e2) {
        let i2 = t2.previousElementSibling;
        for (; i2; ) {
          if (i2.matches(e2))
            return [i2];
          i2 = i2.previousElementSibling;
        }
        return [];
      }, next(t2, e2) {
        let i2 = t2.nextElementSibling;
        for (; i2; ) {
          if (i2.matches(e2))
            return [i2];
          i2 = i2.nextElementSibling;
        }
        return [];
      }, focusableChildren(t2) {
        const e2 = ["a", "button", "input", "textarea", "select", "details", "[tabindex]", '[contenteditable="true"]'].map((t3) => `${t3}:not([tabindex^="-"])`).join(",");
        return this.find(e2, t2).filter((t3) => !l(t3) && a(t3));
      }, getSelectorFromElement(t2) {
        const e2 = B(t2);
        return e2 && z.findOne(e2) ? e2 : null;
      }, getElementFromSelector(t2) {
        const e2 = B(t2);
        return e2 ? z.findOne(e2) : null;
      }, getMultipleElementsFromSelector(t2) {
        const e2 = B(t2);
        return e2 ? z.find(e2) : [];
      } }, R = (t2, e2 = "hide") => {
        const i2 = `click.dismiss${t2.EVENT_KEY}`, n2 = t2.NAME;
        N.on(document, i2, `[data-bs-dismiss="${n2}"]`, function(i3) {
          if (["A", "AREA"].includes(this.tagName) && i3.preventDefault(), l(this))
            return;
          const s2 = z.getElementFromSelector(this) || this.closest(`.${n2}`);
          t2.getOrCreateInstance(s2)[e2]();
        });
      }, q = ".bs.alert", V = `close${q}`, K = `closed${q}`;
      class Q extends W {
        static get NAME() {
          return "alert";
        }
        close() {
          if (N.trigger(this._element, V).defaultPrevented)
            return;
          this._element.classList.remove("show");
          const t2 = this._element.classList.contains("fade");
          this._queueCallback(() => this._destroyElement(), this._element, t2);
        }
        _destroyElement() {
          this._element.remove(), N.trigger(this._element, K), this.dispose();
        }
        static jQueryInterface(t2) {
          return this.each(function() {
            const e2 = Q.getOrCreateInstance(this);
            if ("string" == typeof t2) {
              if (void 0 === e2[t2] || t2.startsWith("_") || "constructor" === t2)
                throw new TypeError(`No method named "${t2}"`);
              e2[t2](this);
            }
          });
        }
      }
      R(Q, "close"), m(Q);
      const X = '[data-bs-toggle="button"]';
      class Y extends W {
        static get NAME() {
          return "button";
        }
        toggle() {
          this._element.setAttribute("aria-pressed", this._element.classList.toggle("active"));
        }
        static jQueryInterface(t2) {
          return this.each(function() {
            const e2 = Y.getOrCreateInstance(this);
            "toggle" === t2 && e2[t2]();
          });
        }
      }
      N.on(document, "click.bs.button.data-api", X, (t2) => {
        t2.preventDefault();
        const e2 = t2.target.closest(X);
        Y.getOrCreateInstance(e2).toggle();
      }), m(Y);
      const U = ".bs.swipe", G = `touchstart${U}`, J = `touchmove${U}`, Z = `touchend${U}`, tt = `pointerdown${U}`, et = `pointerup${U}`, it = { endCallback: null, leftCallback: null, rightCallback: null }, nt = { endCallback: "(function|null)", leftCallback: "(function|null)", rightCallback: "(function|null)" };
      class st extends H {
        constructor(t2, e2) {
          super(), this._element = t2, t2 && st.isSupported() && (this._config = this._getConfig(e2), this._deltaX = 0, this._supportPointerEvents = Boolean(window.PointerEvent), this._initEvents());
        }
        static get Default() {
          return it;
        }
        static get DefaultType() {
          return nt;
        }
        static get NAME() {
          return "swipe";
        }
        dispose() {
          N.off(this._element, U);
        }
        _start(t2) {
          this._supportPointerEvents ? this._eventIsPointerPenTouch(t2) && (this._deltaX = t2.clientX) : this._deltaX = t2.touches[0].clientX;
        }
        _end(t2) {
          this._eventIsPointerPenTouch(t2) && (this._deltaX = t2.clientX - this._deltaX), this._handleSwipe(), g(this._config.endCallback);
        }
        _move(t2) {
          this._deltaX = t2.touches && t2.touches.length > 1 ? 0 : t2.touches[0].clientX - this._deltaX;
        }
        _handleSwipe() {
          const t2 = Math.abs(this._deltaX);
          if (t2 <= 40)
            return;
          const e2 = t2 / this._deltaX;
          this._deltaX = 0, e2 && g(e2 > 0 ? this._config.rightCallback : this._config.leftCallback);
        }
        _initEvents() {
          this._supportPointerEvents ? (N.on(this._element, tt, (t2) => this._start(t2)), N.on(this._element, et, (t2) => this._end(t2)), this._element.classList.add("pointer-event")) : (N.on(this._element, G, (t2) => this._start(t2)), N.on(this._element, J, (t2) => this._move(t2)), N.on(this._element, Z, (t2) => this._end(t2)));
        }
        _eventIsPointerPenTouch(t2) {
          return this._supportPointerEvents && ("pen" === t2.pointerType || "touch" === t2.pointerType);
        }
        static isSupported() {
          return "ontouchstart" in document.documentElement || navigator.maxTouchPoints > 0;
        }
      }
      const ot = ".bs.carousel", rt = ".data-api", at = "next", lt = "prev", ct = "left", ht = "right", dt = `slide${ot}`, ut = `slid${ot}`, ft = `keydown${ot}`, pt = `mouseenter${ot}`, mt = `mouseleave${ot}`, gt = `dragstart${ot}`, _t = `load${ot}${rt}`, bt = `click${ot}${rt}`, vt = "carousel", yt = "active", wt = ".active", At = ".carousel-item", Et = wt + At, Tt = { ArrowLeft: ht, ArrowRight: ct }, Ct = { interval: 5e3, keyboard: true, pause: "hover", ride: false, touch: true, wrap: true }, Ot = { interval: "(number|boolean)", keyboard: "boolean", pause: "(string|boolean)", ride: "(boolean|string)", touch: "boolean", wrap: "boolean" };
      class xt extends W {
        constructor(t2, e2) {
          super(t2, e2), this._interval = null, this._activeElement = null, this._isSliding = false, this.touchTimeout = null, this._swipeHelper = null, this._indicatorsElement = z.findOne(".carousel-indicators", this._element), this._addEventListeners(), this._config.ride === vt && this.cycle();
        }
        static get Default() {
          return Ct;
        }
        static get DefaultType() {
          return Ot;
        }
        static get NAME() {
          return "carousel";
        }
        next() {
          this._slide(at);
        }
        nextWhenVisible() {
          !document.hidden && a(this._element) && this.next();
        }
        prev() {
          this._slide(lt);
        }
        pause() {
          this._isSliding && s(this._element), this._clearInterval();
        }
        cycle() {
          this._clearInterval(), this._updateInterval(), this._interval = setInterval(() => this.nextWhenVisible(), this._config.interval);
        }
        _maybeEnableCycle() {
          this._config.ride && (this._isSliding ? N.one(this._element, ut, () => this.cycle()) : this.cycle());
        }
        to(t2) {
          const e2 = this._getItems();
          if (t2 > e2.length - 1 || t2 < 0)
            return;
          if (this._isSliding)
            return void N.one(this._element, ut, () => this.to(t2));
          const i2 = this._getItemIndex(this._getActive());
          if (i2 === t2)
            return;
          const n2 = t2 > i2 ? at : lt;
          this._slide(n2, e2[t2]);
        }
        dispose() {
          this._swipeHelper && this._swipeHelper.dispose(), super.dispose();
        }
        _configAfterMerge(t2) {
          return t2.defaultInterval = t2.interval, t2;
        }
        _addEventListeners() {
          this._config.keyboard && N.on(this._element, ft, (t2) => this._keydown(t2)), "hover" === this._config.pause && (N.on(this._element, pt, () => this.pause()), N.on(this._element, mt, () => this._maybeEnableCycle())), this._config.touch && st.isSupported() && this._addTouchEventListeners();
        }
        _addTouchEventListeners() {
          for (const t3 of z.find(".carousel-item img", this._element))
            N.on(t3, gt, (t4) => t4.preventDefault());
          const t2 = { leftCallback: () => this._slide(this._directionToOrder(ct)), rightCallback: () => this._slide(this._directionToOrder(ht)), endCallback: () => {
            "hover" === this._config.pause && (this.pause(), this.touchTimeout && clearTimeout(this.touchTimeout), this.touchTimeout = setTimeout(() => this._maybeEnableCycle(), 500 + this._config.interval));
          } };
          this._swipeHelper = new st(this._element, t2);
        }
        _keydown(t2) {
          if (/input|textarea/i.test(t2.target.tagName))
            return;
          const e2 = Tt[t2.key];
          e2 && (t2.preventDefault(), this._slide(this._directionToOrder(e2)));
        }
        _getItemIndex(t2) {
          return this._getItems().indexOf(t2);
        }
        _setActiveIndicatorElement(t2) {
          if (!this._indicatorsElement)
            return;
          const e2 = z.findOne(wt, this._indicatorsElement);
          e2.classList.remove(yt), e2.removeAttribute("aria-current");
          const i2 = z.findOne(`[data-bs-slide-to="${t2}"]`, this._indicatorsElement);
          i2 && (i2.classList.add(yt), i2.setAttribute("aria-current", "true"));
        }
        _updateInterval() {
          const t2 = this._activeElement || this._getActive();
          if (!t2)
            return;
          const e2 = Number.parseInt(t2.getAttribute("data-bs-interval"), 10);
          this._config.interval = e2 || this._config.defaultInterval;
        }
        _slide(t2, e2 = null) {
          if (this._isSliding)
            return;
          const i2 = this._getActive(), n2 = t2 === at, s2 = e2 || b(this._getItems(), i2, n2, this._config.wrap);
          if (s2 === i2)
            return;
          const o2 = this._getItemIndex(s2), r2 = (e3) => N.trigger(this._element, e3, { relatedTarget: s2, direction: this._orderToDirection(t2), from: this._getItemIndex(i2), to: o2 });
          if (r2(dt).defaultPrevented)
            return;
          if (!i2 || !s2)
            return;
          const a2 = Boolean(this._interval);
          this.pause(), this._isSliding = true, this._setActiveIndicatorElement(o2), this._activeElement = s2;
          const l2 = n2 ? "carousel-item-start" : "carousel-item-end", c2 = n2 ? "carousel-item-next" : "carousel-item-prev";
          s2.classList.add(c2), d(s2), i2.classList.add(l2), s2.classList.add(l2), this._queueCallback(() => {
            s2.classList.remove(l2, c2), s2.classList.add(yt), i2.classList.remove(yt, c2, l2), this._isSliding = false, r2(ut);
          }, i2, this._isAnimated()), a2 && this.cycle();
        }
        _isAnimated() {
          return this._element.classList.contains("slide");
        }
        _getActive() {
          return z.findOne(Et, this._element);
        }
        _getItems() {
          return z.find(At, this._element);
        }
        _clearInterval() {
          this._interval && (clearInterval(this._interval), this._interval = null);
        }
        _directionToOrder(t2) {
          return p() ? t2 === ct ? lt : at : t2 === ct ? at : lt;
        }
        _orderToDirection(t2) {
          return p() ? t2 === lt ? ct : ht : t2 === lt ? ht : ct;
        }
        static jQueryInterface(t2) {
          return this.each(function() {
            const e2 = xt.getOrCreateInstance(this, t2);
            if ("number" != typeof t2) {
              if ("string" == typeof t2) {
                if (void 0 === e2[t2] || t2.startsWith("_") || "constructor" === t2)
                  throw new TypeError(`No method named "${t2}"`);
                e2[t2]();
              }
            } else
              e2.to(t2);
          });
        }
      }
      N.on(document, bt, "[data-bs-slide], [data-bs-slide-to]", function(t2) {
        const e2 = z.getElementFromSelector(this);
        if (!e2 || !e2.classList.contains(vt))
          return;
        t2.preventDefault();
        const i2 = xt.getOrCreateInstance(e2), n2 = this.getAttribute("data-bs-slide-to");
        return n2 ? (i2.to(n2), void i2._maybeEnableCycle()) : "next" === F.getDataAttribute(this, "slide") ? (i2.next(), void i2._maybeEnableCycle()) : (i2.prev(), void i2._maybeEnableCycle());
      }), N.on(window, _t, () => {
        const t2 = z.find('[data-bs-ride="carousel"]');
        for (const e2 of t2)
          xt.getOrCreateInstance(e2);
      }), m(xt);
      const kt = ".bs.collapse", Lt = `show${kt}`, St = `shown${kt}`, Dt = `hide${kt}`, $t = `hidden${kt}`, It = `click${kt}.data-api`, Nt = "show", Pt = "collapse", Mt = "collapsing", jt = `:scope .${Pt} .${Pt}`, Ft = '[data-bs-toggle="collapse"]', Ht = { parent: null, toggle: true }, Wt = { parent: "(null|element)", toggle: "boolean" };
      class Bt extends W {
        constructor(t2, e2) {
          super(t2, e2), this._isTransitioning = false, this._triggerArray = [];
          const i2 = z.find(Ft);
          for (const t3 of i2) {
            const e3 = z.getSelectorFromElement(t3), i3 = z.find(e3).filter((t4) => t4 === this._element);
            null !== e3 && i3.length && this._triggerArray.push(t3);
          }
          this._initializeChildren(), this._config.parent || this._addAriaAndCollapsedClass(this._triggerArray, this._isShown()), this._config.toggle && this.toggle();
        }
        static get Default() {
          return Ht;
        }
        static get DefaultType() {
          return Wt;
        }
        static get NAME() {
          return "collapse";
        }
        toggle() {
          this._isShown() ? this.hide() : this.show();
        }
        show() {
          if (this._isTransitioning || this._isShown())
            return;
          let t2 = [];
          if (this._config.parent && (t2 = this._getFirstLevelChildren(".collapse.show, .collapse.collapsing").filter((t3) => t3 !== this._element).map((t3) => Bt.getOrCreateInstance(t3, { toggle: false }))), t2.length && t2[0]._isTransitioning)
            return;
          if (N.trigger(this._element, Lt).defaultPrevented)
            return;
          for (const e3 of t2)
            e3.hide();
          const e2 = this._getDimension();
          this._element.classList.remove(Pt), this._element.classList.add(Mt), this._element.style[e2] = 0, this._addAriaAndCollapsedClass(this._triggerArray, true), this._isTransitioning = true;
          const i2 = `scroll${e2[0].toUpperCase() + e2.slice(1)}`;
          this._queueCallback(() => {
            this._isTransitioning = false, this._element.classList.remove(Mt), this._element.classList.add(Pt, Nt), this._element.style[e2] = "", N.trigger(this._element, St);
          }, this._element, true), this._element.style[e2] = `${this._element[i2]}px`;
        }
        hide() {
          if (this._isTransitioning || !this._isShown())
            return;
          if (N.trigger(this._element, Dt).defaultPrevented)
            return;
          const t2 = this._getDimension();
          this._element.style[t2] = `${this._element.getBoundingClientRect()[t2]}px`, d(this._element), this._element.classList.add(Mt), this._element.classList.remove(Pt, Nt);
          for (const t3 of this._triggerArray) {
            const e2 = z.getElementFromSelector(t3);
            e2 && !this._isShown(e2) && this._addAriaAndCollapsedClass([t3], false);
          }
          this._isTransitioning = true, this._element.style[t2] = "", this._queueCallback(() => {
            this._isTransitioning = false, this._element.classList.remove(Mt), this._element.classList.add(Pt), N.trigger(this._element, $t);
          }, this._element, true);
        }
        _isShown(t2 = this._element) {
          return t2.classList.contains(Nt);
        }
        _configAfterMerge(t2) {
          return t2.toggle = Boolean(t2.toggle), t2.parent = r(t2.parent), t2;
        }
        _getDimension() {
          return this._element.classList.contains("collapse-horizontal") ? "width" : "height";
        }
        _initializeChildren() {
          if (!this._config.parent)
            return;
          const t2 = this._getFirstLevelChildren(Ft);
          for (const e2 of t2) {
            const t3 = z.getElementFromSelector(e2);
            t3 && this._addAriaAndCollapsedClass([e2], this._isShown(t3));
          }
        }
        _getFirstLevelChildren(t2) {
          const e2 = z.find(jt, this._config.parent);
          return z.find(t2, this._config.parent).filter((t3) => !e2.includes(t3));
        }
        _addAriaAndCollapsedClass(t2, e2) {
          if (t2.length)
            for (const i2 of t2)
              i2.classList.toggle("collapsed", !e2), i2.setAttribute("aria-expanded", e2);
        }
        static jQueryInterface(t2) {
          const e2 = {};
          return "string" == typeof t2 && /show|hide/.test(t2) && (e2.toggle = false), this.each(function() {
            const i2 = Bt.getOrCreateInstance(this, e2);
            if ("string" == typeof t2) {
              if (void 0 === i2[t2])
                throw new TypeError(`No method named "${t2}"`);
              i2[t2]();
            }
          });
        }
      }
      N.on(document, It, Ft, function(t2) {
        ("A" === t2.target.tagName || t2.delegateTarget && "A" === t2.delegateTarget.tagName) && t2.preventDefault();
        for (const t3 of z.getMultipleElementsFromSelector(this))
          Bt.getOrCreateInstance(t3, { toggle: false }).toggle();
      }), m(Bt);
      var zt = "top", Rt = "bottom", qt = "right", Vt = "left", Kt = "auto", Qt = [zt, Rt, qt, Vt], Xt = "start", Yt = "end", Ut = "clippingParents", Gt = "viewport", Jt = "popper", Zt = "reference", te = Qt.reduce(function(t2, e2) {
        return t2.concat([e2 + "-" + Xt, e2 + "-" + Yt]);
      }, []), ee = [].concat(Qt, [Kt]).reduce(function(t2, e2) {
        return t2.concat([e2, e2 + "-" + Xt, e2 + "-" + Yt]);
      }, []), ie = "beforeRead", ne = "read", se = "afterRead", oe = "beforeMain", re = "main", ae = "afterMain", le = "beforeWrite", ce = "write", he = "afterWrite", de = [ie, ne, se, oe, re, ae, le, ce, he];
      function ue(t2) {
        return t2 ? (t2.nodeName || "").toLowerCase() : null;
      }
      function fe(t2) {
        if (null == t2)
          return window;
        if ("[object Window]" !== t2.toString()) {
          var e2 = t2.ownerDocument;
          return e2 && e2.defaultView || window;
        }
        return t2;
      }
      function pe(t2) {
        return t2 instanceof fe(t2).Element || t2 instanceof Element;
      }
      function me(t2) {
        return t2 instanceof fe(t2).HTMLElement || t2 instanceof HTMLElement;
      }
      function ge(t2) {
        return "undefined" != typeof ShadowRoot && (t2 instanceof fe(t2).ShadowRoot || t2 instanceof ShadowRoot);
      }
      const _e = { name: "applyStyles", enabled: true, phase: "write", fn: function(t2) {
        var e2 = t2.state;
        Object.keys(e2.elements).forEach(function(t3) {
          var i2 = e2.styles[t3] || {}, n2 = e2.attributes[t3] || {}, s2 = e2.elements[t3];
          me(s2) && ue(s2) && (Object.assign(s2.style, i2), Object.keys(n2).forEach(function(t4) {
            var e3 = n2[t4];
            false === e3 ? s2.removeAttribute(t4) : s2.setAttribute(t4, true === e3 ? "" : e3);
          }));
        });
      }, effect: function(t2) {
        var e2 = t2.state, i2 = { popper: { position: e2.options.strategy, left: "0", top: "0", margin: "0" }, arrow: { position: "absolute" }, reference: {} };
        return Object.assign(e2.elements.popper.style, i2.popper), e2.styles = i2, e2.elements.arrow && Object.assign(e2.elements.arrow.style, i2.arrow), function() {
          Object.keys(e2.elements).forEach(function(t3) {
            var n2 = e2.elements[t3], s2 = e2.attributes[t3] || {}, o2 = Object.keys(e2.styles.hasOwnProperty(t3) ? e2.styles[t3] : i2[t3]).reduce(function(t4, e3) {
              return t4[e3] = "", t4;
            }, {});
            me(n2) && ue(n2) && (Object.assign(n2.style, o2), Object.keys(s2).forEach(function(t4) {
              n2.removeAttribute(t4);
            }));
          });
        };
      }, requires: ["computeStyles"] };
      function be(t2) {
        return t2.split("-")[0];
      }
      var ve = Math.max, ye = Math.min, we = Math.round;
      function Ae() {
        var t2 = navigator.userAgentData;
        return null != t2 && t2.brands && Array.isArray(t2.brands) ? t2.brands.map(function(t3) {
          return t3.brand + "/" + t3.version;
        }).join(" ") : navigator.userAgent;
      }
      function Ee() {
        return !/^((?!chrome|android).)*safari/i.test(Ae());
      }
      function Te(t2, e2, i2) {
        void 0 === e2 && (e2 = false), void 0 === i2 && (i2 = false);
        var n2 = t2.getBoundingClientRect(), s2 = 1, o2 = 1;
        e2 && me(t2) && (s2 = t2.offsetWidth > 0 && we(n2.width) / t2.offsetWidth || 1, o2 = t2.offsetHeight > 0 && we(n2.height) / t2.offsetHeight || 1);
        var r2 = (pe(t2) ? fe(t2) : window).visualViewport, a2 = !Ee() && i2, l2 = (n2.left + (a2 && r2 ? r2.offsetLeft : 0)) / s2, c2 = (n2.top + (a2 && r2 ? r2.offsetTop : 0)) / o2, h2 = n2.width / s2, d2 = n2.height / o2;
        return { width: h2, height: d2, top: c2, right: l2 + h2, bottom: c2 + d2, left: l2, x: l2, y: c2 };
      }
      function Ce(t2) {
        var e2 = Te(t2), i2 = t2.offsetWidth, n2 = t2.offsetHeight;
        return Math.abs(e2.width - i2) <= 1 && (i2 = e2.width), Math.abs(e2.height - n2) <= 1 && (n2 = e2.height), { x: t2.offsetLeft, y: t2.offsetTop, width: i2, height: n2 };
      }
      function Oe(t2, e2) {
        var i2 = e2.getRootNode && e2.getRootNode();
        if (t2.contains(e2))
          return true;
        if (i2 && ge(i2)) {
          var n2 = e2;
          do {
            if (n2 && t2.isSameNode(n2))
              return true;
            n2 = n2.parentNode || n2.host;
          } while (n2);
        }
        return false;
      }
      function xe(t2) {
        return fe(t2).getComputedStyle(t2);
      }
      function ke(t2) {
        return ["table", "td", "th"].indexOf(ue(t2)) >= 0;
      }
      function Le(t2) {
        return ((pe(t2) ? t2.ownerDocument : t2.document) || window.document).documentElement;
      }
      function Se(t2) {
        return "html" === ue(t2) ? t2 : t2.assignedSlot || t2.parentNode || (ge(t2) ? t2.host : null) || Le(t2);
      }
      function De(t2) {
        return me(t2) && "fixed" !== xe(t2).position ? t2.offsetParent : null;
      }
      function $e(t2) {
        for (var e2 = fe(t2), i2 = De(t2); i2 && ke(i2) && "static" === xe(i2).position; )
          i2 = De(i2);
        return i2 && ("html" === ue(i2) || "body" === ue(i2) && "static" === xe(i2).position) ? e2 : i2 || function(t3) {
          var e3 = /firefox/i.test(Ae());
          if (/Trident/i.test(Ae()) && me(t3) && "fixed" === xe(t3).position)
            return null;
          var i3 = Se(t3);
          for (ge(i3) && (i3 = i3.host); me(i3) && ["html", "body"].indexOf(ue(i3)) < 0; ) {
            var n2 = xe(i3);
            if ("none" !== n2.transform || "none" !== n2.perspective || "paint" === n2.contain || -1 !== ["transform", "perspective"].indexOf(n2.willChange) || e3 && "filter" === n2.willChange || e3 && n2.filter && "none" !== n2.filter)
              return i3;
            i3 = i3.parentNode;
          }
          return null;
        }(t2) || e2;
      }
      function Ie(t2) {
        return ["top", "bottom"].indexOf(t2) >= 0 ? "x" : "y";
      }
      function Ne(t2, e2, i2) {
        return ve(t2, ye(e2, i2));
      }
      function Pe(t2) {
        return Object.assign({}, { top: 0, right: 0, bottom: 0, left: 0 }, t2);
      }
      function Me(t2, e2) {
        return e2.reduce(function(e3, i2) {
          return e3[i2] = t2, e3;
        }, {});
      }
      const je = { name: "arrow", enabled: true, phase: "main", fn: function(t2) {
        var e2, i2 = t2.state, n2 = t2.name, s2 = t2.options, o2 = i2.elements.arrow, r2 = i2.modifiersData.popperOffsets, a2 = be(i2.placement), l2 = Ie(a2), c2 = [Vt, qt].indexOf(a2) >= 0 ? "height" : "width";
        if (o2 && r2) {
          var h2 = function(t3, e3) {
            return Pe("number" != typeof (t3 = "function" == typeof t3 ? t3(Object.assign({}, e3.rects, { placement: e3.placement })) : t3) ? t3 : Me(t3, Qt));
          }(s2.padding, i2), d2 = Ce(o2), u2 = "y" === l2 ? zt : Vt, f2 = "y" === l2 ? Rt : qt, p2 = i2.rects.reference[c2] + i2.rects.reference[l2] - r2[l2] - i2.rects.popper[c2], m2 = r2[l2] - i2.rects.reference[l2], g2 = $e(o2), _2 = g2 ? "y" === l2 ? g2.clientHeight || 0 : g2.clientWidth || 0 : 0, b2 = p2 / 2 - m2 / 2, v2 = h2[u2], y2 = _2 - d2[c2] - h2[f2], w2 = _2 / 2 - d2[c2] / 2 + b2, A2 = Ne(v2, w2, y2), E2 = l2;
          i2.modifiersData[n2] = ((e2 = {})[E2] = A2, e2.centerOffset = A2 - w2, e2);
        }
      }, effect: function(t2) {
        var e2 = t2.state, i2 = t2.options.element, n2 = void 0 === i2 ? "[data-popper-arrow]" : i2;
        null != n2 && ("string" != typeof n2 || (n2 = e2.elements.popper.querySelector(n2))) && Oe(e2.elements.popper, n2) && (e2.elements.arrow = n2);
      }, requires: ["popperOffsets"], requiresIfExists: ["preventOverflow"] };
      function Fe(t2) {
        return t2.split("-")[1];
      }
      var He = { top: "auto", right: "auto", bottom: "auto", left: "auto" };
      function We(t2) {
        var e2, i2 = t2.popper, n2 = t2.popperRect, s2 = t2.placement, o2 = t2.variation, r2 = t2.offsets, a2 = t2.position, l2 = t2.gpuAcceleration, c2 = t2.adaptive, h2 = t2.roundOffsets, d2 = t2.isFixed, u2 = r2.x, f2 = void 0 === u2 ? 0 : u2, p2 = r2.y, m2 = void 0 === p2 ? 0 : p2, g2 = "function" == typeof h2 ? h2({ x: f2, y: m2 }) : { x: f2, y: m2 };
        f2 = g2.x, m2 = g2.y;
        var _2 = r2.hasOwnProperty("x"), b2 = r2.hasOwnProperty("y"), v2 = Vt, y2 = zt, w2 = window;
        if (c2) {
          var A2 = $e(i2), E2 = "clientHeight", T2 = "clientWidth";
          A2 === fe(i2) && "static" !== xe(A2 = Le(i2)).position && "absolute" === a2 && (E2 = "scrollHeight", T2 = "scrollWidth"), (s2 === zt || (s2 === Vt || s2 === qt) && o2 === Yt) && (y2 = Rt, m2 -= (d2 && A2 === w2 && w2.visualViewport ? w2.visualViewport.height : A2[E2]) - n2.height, m2 *= l2 ? 1 : -1), s2 !== Vt && (s2 !== zt && s2 !== Rt || o2 !== Yt) || (v2 = qt, f2 -= (d2 && A2 === w2 && w2.visualViewport ? w2.visualViewport.width : A2[T2]) - n2.width, f2 *= l2 ? 1 : -1);
        }
        var C2, O2 = Object.assign({ position: a2 }, c2 && He), x2 = true === h2 ? function(t3, e3) {
          var i3 = t3.x, n3 = t3.y, s3 = e3.devicePixelRatio || 1;
          return { x: we(i3 * s3) / s3 || 0, y: we(n3 * s3) / s3 || 0 };
        }({ x: f2, y: m2 }, fe(i2)) : { x: f2, y: m2 };
        return f2 = x2.x, m2 = x2.y, l2 ? Object.assign({}, O2, ((C2 = {})[y2] = b2 ? "0" : "", C2[v2] = _2 ? "0" : "", C2.transform = (w2.devicePixelRatio || 1) <= 1 ? "translate(" + f2 + "px, " + m2 + "px)" : "translate3d(" + f2 + "px, " + m2 + "px, 0)", C2)) : Object.assign({}, O2, ((e2 = {})[y2] = b2 ? m2 + "px" : "", e2[v2] = _2 ? f2 + "px" : "", e2.transform = "", e2));
      }
      const Be = { name: "computeStyles", enabled: true, phase: "beforeWrite", fn: function(t2) {
        var e2 = t2.state, i2 = t2.options, n2 = i2.gpuAcceleration, s2 = void 0 === n2 || n2, o2 = i2.adaptive, r2 = void 0 === o2 || o2, a2 = i2.roundOffsets, l2 = void 0 === a2 || a2, c2 = { placement: be(e2.placement), variation: Fe(e2.placement), popper: e2.elements.popper, popperRect: e2.rects.popper, gpuAcceleration: s2, isFixed: "fixed" === e2.options.strategy };
        null != e2.modifiersData.popperOffsets && (e2.styles.popper = Object.assign({}, e2.styles.popper, We(Object.assign({}, c2, { offsets: e2.modifiersData.popperOffsets, position: e2.options.strategy, adaptive: r2, roundOffsets: l2 })))), null != e2.modifiersData.arrow && (e2.styles.arrow = Object.assign({}, e2.styles.arrow, We(Object.assign({}, c2, { offsets: e2.modifiersData.arrow, position: "absolute", adaptive: false, roundOffsets: l2 })))), e2.attributes.popper = Object.assign({}, e2.attributes.popper, { "data-popper-placement": e2.placement });
      }, data: {} };
      var ze = { passive: true };
      const Re = { name: "eventListeners", enabled: true, phase: "write", fn: function() {
      }, effect: function(t2) {
        var e2 = t2.state, i2 = t2.instance, n2 = t2.options, s2 = n2.scroll, o2 = void 0 === s2 || s2, r2 = n2.resize, a2 = void 0 === r2 || r2, l2 = fe(e2.elements.popper), c2 = [].concat(e2.scrollParents.reference, e2.scrollParents.popper);
        return o2 && c2.forEach(function(t3) {
          t3.addEventListener("scroll", i2.update, ze);
        }), a2 && l2.addEventListener("resize", i2.update, ze), function() {
          o2 && c2.forEach(function(t3) {
            t3.removeEventListener("scroll", i2.update, ze);
          }), a2 && l2.removeEventListener("resize", i2.update, ze);
        };
      }, data: {} };
      var qe = { left: "right", right: "left", bottom: "top", top: "bottom" };
      function Ve(t2) {
        return t2.replace(/left|right|bottom|top/g, function(t3) {
          return qe[t3];
        });
      }
      var Ke = { start: "end", end: "start" };
      function Qe(t2) {
        return t2.replace(/start|end/g, function(t3) {
          return Ke[t3];
        });
      }
      function Xe(t2) {
        var e2 = fe(t2);
        return { scrollLeft: e2.pageXOffset, scrollTop: e2.pageYOffset };
      }
      function Ye(t2) {
        return Te(Le(t2)).left + Xe(t2).scrollLeft;
      }
      function Ue(t2) {
        var e2 = xe(t2), i2 = e2.overflow, n2 = e2.overflowX, s2 = e2.overflowY;
        return /auto|scroll|overlay|hidden/.test(i2 + s2 + n2);
      }
      function Ge(t2) {
        return ["html", "body", "#document"].indexOf(ue(t2)) >= 0 ? t2.ownerDocument.body : me(t2) && Ue(t2) ? t2 : Ge(Se(t2));
      }
      function Je(t2, e2) {
        var i2;
        void 0 === e2 && (e2 = []);
        var n2 = Ge(t2), s2 = n2 === (null == (i2 = t2.ownerDocument) ? void 0 : i2.body), o2 = fe(n2), r2 = s2 ? [o2].concat(o2.visualViewport || [], Ue(n2) ? n2 : []) : n2, a2 = e2.concat(r2);
        return s2 ? a2 : a2.concat(Je(Se(r2)));
      }
      function Ze(t2) {
        return Object.assign({}, t2, { left: t2.x, top: t2.y, right: t2.x + t2.width, bottom: t2.y + t2.height });
      }
      function ti(t2, e2, i2) {
        return e2 === Gt ? Ze(function(t3, e3) {
          var i3 = fe(t3), n2 = Le(t3), s2 = i3.visualViewport, o2 = n2.clientWidth, r2 = n2.clientHeight, a2 = 0, l2 = 0;
          if (s2) {
            o2 = s2.width, r2 = s2.height;
            var c2 = Ee();
            (c2 || !c2 && "fixed" === e3) && (a2 = s2.offsetLeft, l2 = s2.offsetTop);
          }
          return { width: o2, height: r2, x: a2 + Ye(t3), y: l2 };
        }(t2, i2)) : pe(e2) ? function(t3, e3) {
          var i3 = Te(t3, false, "fixed" === e3);
          return i3.top = i3.top + t3.clientTop, i3.left = i3.left + t3.clientLeft, i3.bottom = i3.top + t3.clientHeight, i3.right = i3.left + t3.clientWidth, i3.width = t3.clientWidth, i3.height = t3.clientHeight, i3.x = i3.left, i3.y = i3.top, i3;
        }(e2, i2) : Ze(function(t3) {
          var e3, i3 = Le(t3), n2 = Xe(t3), s2 = null == (e3 = t3.ownerDocument) ? void 0 : e3.body, o2 = ve(i3.scrollWidth, i3.clientWidth, s2 ? s2.scrollWidth : 0, s2 ? s2.clientWidth : 0), r2 = ve(i3.scrollHeight, i3.clientHeight, s2 ? s2.scrollHeight : 0, s2 ? s2.clientHeight : 0), a2 = -n2.scrollLeft + Ye(t3), l2 = -n2.scrollTop;
          return "rtl" === xe(s2 || i3).direction && (a2 += ve(i3.clientWidth, s2 ? s2.clientWidth : 0) - o2), { width: o2, height: r2, x: a2, y: l2 };
        }(Le(t2)));
      }
      function ei(t2) {
        var e2, i2 = t2.reference, n2 = t2.element, s2 = t2.placement, o2 = s2 ? be(s2) : null, r2 = s2 ? Fe(s2) : null, a2 = i2.x + i2.width / 2 - n2.width / 2, l2 = i2.y + i2.height / 2 - n2.height / 2;
        switch (o2) {
          case zt:
            e2 = { x: a2, y: i2.y - n2.height };
            break;
          case Rt:
            e2 = { x: a2, y: i2.y + i2.height };
            break;
          case qt:
            e2 = { x: i2.x + i2.width, y: l2 };
            break;
          case Vt:
            e2 = { x: i2.x - n2.width, y: l2 };
            break;
          default:
            e2 = { x: i2.x, y: i2.y };
        }
        var c2 = o2 ? Ie(o2) : null;
        if (null != c2) {
          var h2 = "y" === c2 ? "height" : "width";
          switch (r2) {
            case Xt:
              e2[c2] = e2[c2] - (i2[h2] / 2 - n2[h2] / 2);
              break;
            case Yt:
              e2[c2] = e2[c2] + (i2[h2] / 2 - n2[h2] / 2);
          }
        }
        return e2;
      }
      function ii(t2, e2) {
        void 0 === e2 && (e2 = {});
        var i2 = e2, n2 = i2.placement, s2 = void 0 === n2 ? t2.placement : n2, o2 = i2.strategy, r2 = void 0 === o2 ? t2.strategy : o2, a2 = i2.boundary, l2 = void 0 === a2 ? Ut : a2, c2 = i2.rootBoundary, h2 = void 0 === c2 ? Gt : c2, d2 = i2.elementContext, u2 = void 0 === d2 ? Jt : d2, f2 = i2.altBoundary, p2 = void 0 !== f2 && f2, m2 = i2.padding, g2 = void 0 === m2 ? 0 : m2, _2 = Pe("number" != typeof g2 ? g2 : Me(g2, Qt)), b2 = u2 === Jt ? Zt : Jt, v2 = t2.rects.popper, y2 = t2.elements[p2 ? b2 : u2], w2 = function(t3, e3, i3, n3) {
          var s3 = "clippingParents" === e3 ? function(t4) {
            var e4 = Je(Se(t4)), i4 = ["absolute", "fixed"].indexOf(xe(t4).position) >= 0 && me(t4) ? $e(t4) : t4;
            return pe(i4) ? e4.filter(function(t5) {
              return pe(t5) && Oe(t5, i4) && "body" !== ue(t5);
            }) : [];
          }(t3) : [].concat(e3), o3 = [].concat(s3, [i3]), r3 = o3[0], a3 = o3.reduce(function(e4, i4) {
            var s4 = ti(t3, i4, n3);
            return e4.top = ve(s4.top, e4.top), e4.right = ye(s4.right, e4.right), e4.bottom = ye(s4.bottom, e4.bottom), e4.left = ve(s4.left, e4.left), e4;
          }, ti(t3, r3, n3));
          return a3.width = a3.right - a3.left, a3.height = a3.bottom - a3.top, a3.x = a3.left, a3.y = a3.top, a3;
        }(pe(y2) ? y2 : y2.contextElement || Le(t2.elements.popper), l2, h2, r2), A2 = Te(t2.elements.reference), E2 = ei({ reference: A2, element: v2, strategy: "absolute", placement: s2 }), T2 = Ze(Object.assign({}, v2, E2)), C2 = u2 === Jt ? T2 : A2, O2 = { top: w2.top - C2.top + _2.top, bottom: C2.bottom - w2.bottom + _2.bottom, left: w2.left - C2.left + _2.left, right: C2.right - w2.right + _2.right }, x2 = t2.modifiersData.offset;
        if (u2 === Jt && x2) {
          var k2 = x2[s2];
          Object.keys(O2).forEach(function(t3) {
            var e3 = [qt, Rt].indexOf(t3) >= 0 ? 1 : -1, i3 = [zt, Rt].indexOf(t3) >= 0 ? "y" : "x";
            O2[t3] += k2[i3] * e3;
          });
        }
        return O2;
      }
      function ni(t2, e2) {
        void 0 === e2 && (e2 = {});
        var i2 = e2, n2 = i2.placement, s2 = i2.boundary, o2 = i2.rootBoundary, r2 = i2.padding, a2 = i2.flipVariations, l2 = i2.allowedAutoPlacements, c2 = void 0 === l2 ? ee : l2, h2 = Fe(n2), d2 = h2 ? a2 ? te : te.filter(function(t3) {
          return Fe(t3) === h2;
        }) : Qt, u2 = d2.filter(function(t3) {
          return c2.indexOf(t3) >= 0;
        });
        0 === u2.length && (u2 = d2);
        var f2 = u2.reduce(function(e3, i3) {
          return e3[i3] = ii(t2, { placement: i3, boundary: s2, rootBoundary: o2, padding: r2 })[be(i3)], e3;
        }, {});
        return Object.keys(f2).sort(function(t3, e3) {
          return f2[t3] - f2[e3];
        });
      }
      const si = { name: "flip", enabled: true, phase: "main", fn: function(t2) {
        var e2 = t2.state, i2 = t2.options, n2 = t2.name;
        if (!e2.modifiersData[n2]._skip) {
          for (var s2 = i2.mainAxis, o2 = void 0 === s2 || s2, r2 = i2.altAxis, a2 = void 0 === r2 || r2, l2 = i2.fallbackPlacements, c2 = i2.padding, h2 = i2.boundary, d2 = i2.rootBoundary, u2 = i2.altBoundary, f2 = i2.flipVariations, p2 = void 0 === f2 || f2, m2 = i2.allowedAutoPlacements, g2 = e2.options.placement, _2 = be(g2), b2 = l2 || (_2 !== g2 && p2 ? function(t3) {
            if (be(t3) === Kt)
              return [];
            var e3 = Ve(t3);
            return [Qe(t3), e3, Qe(e3)];
          }(g2) : [Ve(g2)]), v2 = [g2].concat(b2).reduce(function(t3, i3) {
            return t3.concat(be(i3) === Kt ? ni(e2, { placement: i3, boundary: h2, rootBoundary: d2, padding: c2, flipVariations: p2, allowedAutoPlacements: m2 }) : i3);
          }, []), y2 = e2.rects.reference, w2 = e2.rects.popper, A2 = /* @__PURE__ */ new Map(), E2 = true, T2 = v2[0], C2 = 0; C2 < v2.length; C2++) {
            var O2 = v2[C2], x2 = be(O2), k2 = Fe(O2) === Xt, L2 = [zt, Rt].indexOf(x2) >= 0, S2 = L2 ? "width" : "height", D2 = ii(e2, { placement: O2, boundary: h2, rootBoundary: d2, altBoundary: u2, padding: c2 }), $2 = L2 ? k2 ? qt : Vt : k2 ? Rt : zt;
            y2[S2] > w2[S2] && ($2 = Ve($2));
            var I2 = Ve($2), N2 = [];
            if (o2 && N2.push(D2[x2] <= 0), a2 && N2.push(D2[$2] <= 0, D2[I2] <= 0), N2.every(function(t3) {
              return t3;
            })) {
              T2 = O2, E2 = false;
              break;
            }
            A2.set(O2, N2);
          }
          if (E2)
            for (var P2 = function(t3) {
              var e3 = v2.find(function(e4) {
                var i3 = A2.get(e4);
                if (i3)
                  return i3.slice(0, t3).every(function(t4) {
                    return t4;
                  });
              });
              if (e3)
                return T2 = e3, "break";
            }, M2 = p2 ? 3 : 1; M2 > 0 && "break" !== P2(M2); M2--)
              ;
          e2.placement !== T2 && (e2.modifiersData[n2]._skip = true, e2.placement = T2, e2.reset = true);
        }
      }, requiresIfExists: ["offset"], data: { _skip: false } };
      function oi(t2, e2, i2) {
        return void 0 === i2 && (i2 = { x: 0, y: 0 }), { top: t2.top - e2.height - i2.y, right: t2.right - e2.width + i2.x, bottom: t2.bottom - e2.height + i2.y, left: t2.left - e2.width - i2.x };
      }
      function ri(t2) {
        return [zt, qt, Rt, Vt].some(function(e2) {
          return t2[e2] >= 0;
        });
      }
      const ai = { name: "hide", enabled: true, phase: "main", requiresIfExists: ["preventOverflow"], fn: function(t2) {
        var e2 = t2.state, i2 = t2.name, n2 = e2.rects.reference, s2 = e2.rects.popper, o2 = e2.modifiersData.preventOverflow, r2 = ii(e2, { elementContext: "reference" }), a2 = ii(e2, { altBoundary: true }), l2 = oi(r2, n2), c2 = oi(a2, s2, o2), h2 = ri(l2), d2 = ri(c2);
        e2.modifiersData[i2] = { referenceClippingOffsets: l2, popperEscapeOffsets: c2, isReferenceHidden: h2, hasPopperEscaped: d2 }, e2.attributes.popper = Object.assign({}, e2.attributes.popper, { "data-popper-reference-hidden": h2, "data-popper-escaped": d2 });
      } }, li = { name: "offset", enabled: true, phase: "main", requires: ["popperOffsets"], fn: function(t2) {
        var e2 = t2.state, i2 = t2.options, n2 = t2.name, s2 = i2.offset, o2 = void 0 === s2 ? [0, 0] : s2, r2 = ee.reduce(function(t3, i3) {
          return t3[i3] = function(t4, e3, i4) {
            var n3 = be(t4), s3 = [Vt, zt].indexOf(n3) >= 0 ? -1 : 1, o3 = "function" == typeof i4 ? i4(Object.assign({}, e3, { placement: t4 })) : i4, r3 = o3[0], a3 = o3[1];
            return r3 = r3 || 0, a3 = (a3 || 0) * s3, [Vt, qt].indexOf(n3) >= 0 ? { x: a3, y: r3 } : { x: r3, y: a3 };
          }(i3, e2.rects, o2), t3;
        }, {}), a2 = r2[e2.placement], l2 = a2.x, c2 = a2.y;
        null != e2.modifiersData.popperOffsets && (e2.modifiersData.popperOffsets.x += l2, e2.modifiersData.popperOffsets.y += c2), e2.modifiersData[n2] = r2;
      } }, ci = { name: "popperOffsets", enabled: true, phase: "read", fn: function(t2) {
        var e2 = t2.state, i2 = t2.name;
        e2.modifiersData[i2] = ei({ reference: e2.rects.reference, element: e2.rects.popper, strategy: "absolute", placement: e2.placement });
      }, data: {} }, hi = { name: "preventOverflow", enabled: true, phase: "main", fn: function(t2) {
        var e2 = t2.state, i2 = t2.options, n2 = t2.name, s2 = i2.mainAxis, o2 = void 0 === s2 || s2, r2 = i2.altAxis, a2 = void 0 !== r2 && r2, l2 = i2.boundary, c2 = i2.rootBoundary, h2 = i2.altBoundary, d2 = i2.padding, u2 = i2.tether, f2 = void 0 === u2 || u2, p2 = i2.tetherOffset, m2 = void 0 === p2 ? 0 : p2, g2 = ii(e2, { boundary: l2, rootBoundary: c2, padding: d2, altBoundary: h2 }), _2 = be(e2.placement), b2 = Fe(e2.placement), v2 = !b2, y2 = Ie(_2), w2 = "x" === y2 ? "y" : "x", A2 = e2.modifiersData.popperOffsets, E2 = e2.rects.reference, T2 = e2.rects.popper, C2 = "function" == typeof m2 ? m2(Object.assign({}, e2.rects, { placement: e2.placement })) : m2, O2 = "number" == typeof C2 ? { mainAxis: C2, altAxis: C2 } : Object.assign({ mainAxis: 0, altAxis: 0 }, C2), x2 = e2.modifiersData.offset ? e2.modifiersData.offset[e2.placement] : null, k2 = { x: 0, y: 0 };
        if (A2) {
          if (o2) {
            var L2, S2 = "y" === y2 ? zt : Vt, D2 = "y" === y2 ? Rt : qt, $2 = "y" === y2 ? "height" : "width", I2 = A2[y2], N2 = I2 + g2[S2], P2 = I2 - g2[D2], M2 = f2 ? -T2[$2] / 2 : 0, j2 = b2 === Xt ? E2[$2] : T2[$2], F2 = b2 === Xt ? -T2[$2] : -E2[$2], H2 = e2.elements.arrow, W2 = f2 && H2 ? Ce(H2) : { width: 0, height: 0 }, B2 = e2.modifiersData["arrow#persistent"] ? e2.modifiersData["arrow#persistent"].padding : { top: 0, right: 0, bottom: 0, left: 0 }, z2 = B2[S2], R2 = B2[D2], q2 = Ne(0, E2[$2], W2[$2]), V2 = v2 ? E2[$2] / 2 - M2 - q2 - z2 - O2.mainAxis : j2 - q2 - z2 - O2.mainAxis, K2 = v2 ? -E2[$2] / 2 + M2 + q2 + R2 + O2.mainAxis : F2 + q2 + R2 + O2.mainAxis, Q2 = e2.elements.arrow && $e(e2.elements.arrow), X2 = Q2 ? "y" === y2 ? Q2.clientTop || 0 : Q2.clientLeft || 0 : 0, Y2 = null != (L2 = null == x2 ? void 0 : x2[y2]) ? L2 : 0, U2 = I2 + K2 - Y2, G2 = Ne(f2 ? ye(N2, I2 + V2 - Y2 - X2) : N2, I2, f2 ? ve(P2, U2) : P2);
            A2[y2] = G2, k2[y2] = G2 - I2;
          }
          if (a2) {
            var J2, Z2 = "x" === y2 ? zt : Vt, tt2 = "x" === y2 ? Rt : qt, et2 = A2[w2], it2 = "y" === w2 ? "height" : "width", nt2 = et2 + g2[Z2], st2 = et2 - g2[tt2], ot2 = -1 !== [zt, Vt].indexOf(_2), rt2 = null != (J2 = null == x2 ? void 0 : x2[w2]) ? J2 : 0, at2 = ot2 ? nt2 : et2 - E2[it2] - T2[it2] - rt2 + O2.altAxis, lt2 = ot2 ? et2 + E2[it2] + T2[it2] - rt2 - O2.altAxis : st2, ct2 = f2 && ot2 ? function(t3, e3, i3) {
              var n3 = Ne(t3, e3, i3);
              return n3 > i3 ? i3 : n3;
            }(at2, et2, lt2) : Ne(f2 ? at2 : nt2, et2, f2 ? lt2 : st2);
            A2[w2] = ct2, k2[w2] = ct2 - et2;
          }
          e2.modifiersData[n2] = k2;
        }
      }, requiresIfExists: ["offset"] };
      function di(t2, e2, i2) {
        void 0 === i2 && (i2 = false);
        var n2, s2, o2 = me(e2), r2 = me(e2) && function(t3) {
          var e3 = t3.getBoundingClientRect(), i3 = we(e3.width) / t3.offsetWidth || 1, n3 = we(e3.height) / t3.offsetHeight || 1;
          return 1 !== i3 || 1 !== n3;
        }(e2), a2 = Le(e2), l2 = Te(t2, r2, i2), c2 = { scrollLeft: 0, scrollTop: 0 }, h2 = { x: 0, y: 0 };
        return (o2 || !o2 && !i2) && (("body" !== ue(e2) || Ue(a2)) && (c2 = (n2 = e2) !== fe(n2) && me(n2) ? { scrollLeft: (s2 = n2).scrollLeft, scrollTop: s2.scrollTop } : Xe(n2)), me(e2) ? ((h2 = Te(e2, true)).x += e2.clientLeft, h2.y += e2.clientTop) : a2 && (h2.x = Ye(a2))), { x: l2.left + c2.scrollLeft - h2.x, y: l2.top + c2.scrollTop - h2.y, width: l2.width, height: l2.height };
      }
      function ui(t2) {
        var e2 = /* @__PURE__ */ new Map(), i2 = /* @__PURE__ */ new Set(), n2 = [];
        function s2(t3) {
          i2.add(t3.name), [].concat(t3.requires || [], t3.requiresIfExists || []).forEach(function(t4) {
            if (!i2.has(t4)) {
              var n3 = e2.get(t4);
              n3 && s2(n3);
            }
          }), n2.push(t3);
        }
        return t2.forEach(function(t3) {
          e2.set(t3.name, t3);
        }), t2.forEach(function(t3) {
          i2.has(t3.name) || s2(t3);
        }), n2;
      }
      var fi = { placement: "bottom", modifiers: [], strategy: "absolute" };
      function pi() {
        for (var t2 = arguments.length, e2 = new Array(t2), i2 = 0; i2 < t2; i2++)
          e2[i2] = arguments[i2];
        return !e2.some(function(t3) {
          return !(t3 && "function" == typeof t3.getBoundingClientRect);
        });
      }
      function mi(t2) {
        void 0 === t2 && (t2 = {});
        var e2 = t2, i2 = e2.defaultModifiers, n2 = void 0 === i2 ? [] : i2, s2 = e2.defaultOptions, o2 = void 0 === s2 ? fi : s2;
        return function(t3, e3, i3) {
          void 0 === i3 && (i3 = o2);
          var s3, r2, a2 = { placement: "bottom", orderedModifiers: [], options: Object.assign({}, fi, o2), modifiersData: {}, elements: { reference: t3, popper: e3 }, attributes: {}, styles: {} }, l2 = [], c2 = false, h2 = { state: a2, setOptions: function(i4) {
            var s4 = "function" == typeof i4 ? i4(a2.options) : i4;
            d2(), a2.options = Object.assign({}, o2, a2.options, s4), a2.scrollParents = { reference: pe(t3) ? Je(t3) : t3.contextElement ? Je(t3.contextElement) : [], popper: Je(e3) };
            var r3, c3, u2 = function(t4) {
              var e4 = ui(t4);
              return de.reduce(function(t5, i5) {
                return t5.concat(e4.filter(function(t6) {
                  return t6.phase === i5;
                }));
              }, []);
            }((r3 = [].concat(n2, a2.options.modifiers), c3 = r3.reduce(function(t4, e4) {
              var i5 = t4[e4.name];
              return t4[e4.name] = i5 ? Object.assign({}, i5, e4, { options: Object.assign({}, i5.options, e4.options), data: Object.assign({}, i5.data, e4.data) }) : e4, t4;
            }, {}), Object.keys(c3).map(function(t4) {
              return c3[t4];
            })));
            return a2.orderedModifiers = u2.filter(function(t4) {
              return t4.enabled;
            }), a2.orderedModifiers.forEach(function(t4) {
              var e4 = t4.name, i5 = t4.options, n3 = void 0 === i5 ? {} : i5, s5 = t4.effect;
              if ("function" == typeof s5) {
                var o3 = s5({ state: a2, name: e4, instance: h2, options: n3 });
                l2.push(o3 || function() {
                });
              }
            }), h2.update();
          }, forceUpdate: function() {
            if (!c2) {
              var t4 = a2.elements, e4 = t4.reference, i4 = t4.popper;
              if (pi(e4, i4)) {
                a2.rects = { reference: di(e4, $e(i4), "fixed" === a2.options.strategy), popper: Ce(i4) }, a2.reset = false, a2.placement = a2.options.placement, a2.orderedModifiers.forEach(function(t5) {
                  return a2.modifiersData[t5.name] = Object.assign({}, t5.data);
                });
                for (var n3 = 0; n3 < a2.orderedModifiers.length; n3++)
                  if (true !== a2.reset) {
                    var s4 = a2.orderedModifiers[n3], o3 = s4.fn, r3 = s4.options, l3 = void 0 === r3 ? {} : r3, d3 = s4.name;
                    "function" == typeof o3 && (a2 = o3({ state: a2, options: l3, name: d3, instance: h2 }) || a2);
                  } else
                    a2.reset = false, n3 = -1;
              }
            }
          }, update: (s3 = function() {
            return new Promise(function(t4) {
              h2.forceUpdate(), t4(a2);
            });
          }, function() {
            return r2 || (r2 = new Promise(function(t4) {
              Promise.resolve().then(function() {
                r2 = void 0, t4(s3());
              });
            })), r2;
          }), destroy: function() {
            d2(), c2 = true;
          } };
          if (!pi(t3, e3))
            return h2;
          function d2() {
            l2.forEach(function(t4) {
              return t4();
            }), l2 = [];
          }
          return h2.setOptions(i3).then(function(t4) {
            !c2 && i3.onFirstUpdate && i3.onFirstUpdate(t4);
          }), h2;
        };
      }
      var gi = mi(), _i = mi({ defaultModifiers: [Re, ci, Be, _e] }), bi = mi({ defaultModifiers: [Re, ci, Be, _e, li, si, hi, je, ai] });
      const vi = Object.freeze(Object.defineProperty({ __proto__: null, afterMain: ae, afterRead: se, afterWrite: he, applyStyles: _e, arrow: je, auto: Kt, basePlacements: Qt, beforeMain: oe, beforeRead: ie, beforeWrite: le, bottom: Rt, clippingParents: Ut, computeStyles: Be, createPopper: bi, createPopperBase: gi, createPopperLite: _i, detectOverflow: ii, end: Yt, eventListeners: Re, flip: si, hide: ai, left: Vt, main: re, modifierPhases: de, offset: li, placements: ee, popper: Jt, popperGenerator: mi, popperOffsets: ci, preventOverflow: hi, read: ne, reference: Zt, right: qt, start: Xt, top: zt, variationPlacements: te, viewport: Gt, write: ce }, Symbol.toStringTag, { value: "Module" })), yi = "dropdown", wi = ".bs.dropdown", Ai = ".data-api", Ei = "ArrowUp", Ti = "ArrowDown", Ci = `hide${wi}`, Oi = `hidden${wi}`, xi = `show${wi}`, ki = `shown${wi}`, Li = `click${wi}${Ai}`, Si = `keydown${wi}${Ai}`, Di = `keyup${wi}${Ai}`, $i = "show", Ii = '[data-bs-toggle="dropdown"]:not(.disabled):not(:disabled)', Ni = `${Ii}.${$i}`, Pi = ".dropdown-menu", Mi = p() ? "top-end" : "top-start", ji = p() ? "top-start" : "top-end", Fi = p() ? "bottom-end" : "bottom-start", Hi = p() ? "bottom-start" : "bottom-end", Wi = p() ? "left-start" : "right-start", Bi = p() ? "right-start" : "left-start", zi = { autoClose: true, boundary: "clippingParents", display: "dynamic", offset: [0, 2], popperConfig: null, reference: "toggle" }, Ri = { autoClose: "(boolean|string)", boundary: "(string|element)", display: "string", offset: "(array|string|function)", popperConfig: "(null|object|function)", reference: "(string|element|object)" };
      class qi extends W {
        constructor(t2, e2) {
          super(t2, e2), this._popper = null, this._parent = this._element.parentNode, this._menu = z.next(this._element, Pi)[0] || z.prev(this._element, Pi)[0] || z.findOne(Pi, this._parent), this._inNavbar = this._detectNavbar();
        }
        static get Default() {
          return zi;
        }
        static get DefaultType() {
          return Ri;
        }
        static get NAME() {
          return yi;
        }
        toggle() {
          return this._isShown() ? this.hide() : this.show();
        }
        show() {
          if (l(this._element) || this._isShown())
            return;
          const t2 = { relatedTarget: this._element };
          if (!N.trigger(this._element, xi, t2).defaultPrevented) {
            if (this._createPopper(), "ontouchstart" in document.documentElement && !this._parent.closest(".navbar-nav"))
              for (const t3 of [].concat(...document.body.children))
                N.on(t3, "mouseover", h);
            this._element.focus(), this._element.setAttribute("aria-expanded", true), this._menu.classList.add($i), this._element.classList.add($i), N.trigger(this._element, ki, t2);
          }
        }
        hide() {
          if (l(this._element) || !this._isShown())
            return;
          const t2 = { relatedTarget: this._element };
          this._completeHide(t2);
        }
        dispose() {
          this._popper && this._popper.destroy(), super.dispose();
        }
        update() {
          this._inNavbar = this._detectNavbar(), this._popper && this._popper.update();
        }
        _completeHide(t2) {
          if (!N.trigger(this._element, Ci, t2).defaultPrevented) {
            if ("ontouchstart" in document.documentElement)
              for (const t3 of [].concat(...document.body.children))
                N.off(t3, "mouseover", h);
            this._popper && this._popper.destroy(), this._menu.classList.remove($i), this._element.classList.remove($i), this._element.setAttribute("aria-expanded", "false"), F.removeDataAttribute(this._menu, "popper"), N.trigger(this._element, Oi, t2);
          }
        }
        _getConfig(t2) {
          if ("object" == typeof (t2 = super._getConfig(t2)).reference && !o(t2.reference) && "function" != typeof t2.reference.getBoundingClientRect)
            throw new TypeError(`${yi.toUpperCase()}: Option "reference" provided type "object" without a required "getBoundingClientRect" method.`);
          return t2;
        }
        _createPopper() {
          if (void 0 === vi)
            throw new TypeError("Bootstrap's dropdowns require Popper (https://popper.js.org)");
          let t2 = this._element;
          "parent" === this._config.reference ? t2 = this._parent : o(this._config.reference) ? t2 = r(this._config.reference) : "object" == typeof this._config.reference && (t2 = this._config.reference);
          const e2 = this._getPopperConfig();
          this._popper = bi(t2, this._menu, e2);
        }
        _isShown() {
          return this._menu.classList.contains($i);
        }
        _getPlacement() {
          const t2 = this._parent;
          if (t2.classList.contains("dropend"))
            return Wi;
          if (t2.classList.contains("dropstart"))
            return Bi;
          if (t2.classList.contains("dropup-center"))
            return "top";
          if (t2.classList.contains("dropdown-center"))
            return "bottom";
          const e2 = "end" === getComputedStyle(this._menu).getPropertyValue("--bs-position").trim();
          return t2.classList.contains("dropup") ? e2 ? ji : Mi : e2 ? Hi : Fi;
        }
        _detectNavbar() {
          return null !== this._element.closest(".navbar");
        }
        _getOffset() {
          const { offset: t2 } = this._config;
          return "string" == typeof t2 ? t2.split(",").map((t3) => Number.parseInt(t3, 10)) : "function" == typeof t2 ? (e2) => t2(e2, this._element) : t2;
        }
        _getPopperConfig() {
          const t2 = { placement: this._getPlacement(), modifiers: [{ name: "preventOverflow", options: { boundary: this._config.boundary } }, { name: "offset", options: { offset: this._getOffset() } }] };
          return (this._inNavbar || "static" === this._config.display) && (F.setDataAttribute(this._menu, "popper", "static"), t2.modifiers = [{ name: "applyStyles", enabled: false }]), { ...t2, ...g(this._config.popperConfig, [t2]) };
        }
        _selectMenuItem({ key: t2, target: e2 }) {
          const i2 = z.find(".dropdown-menu .dropdown-item:not(.disabled):not(:disabled)", this._menu).filter((t3) => a(t3));
          i2.length && b(i2, e2, t2 === Ti, !i2.includes(e2)).focus();
        }
        static jQueryInterface(t2) {
          return this.each(function() {
            const e2 = qi.getOrCreateInstance(this, t2);
            if ("string" == typeof t2) {
              if (void 0 === e2[t2])
                throw new TypeError(`No method named "${t2}"`);
              e2[t2]();
            }
          });
        }
        static clearMenus(t2) {
          if (2 === t2.button || "keyup" === t2.type && "Tab" !== t2.key)
            return;
          const e2 = z.find(Ni);
          for (const i2 of e2) {
            const e3 = qi.getInstance(i2);
            if (!e3 || false === e3._config.autoClose)
              continue;
            const n2 = t2.composedPath(), s2 = n2.includes(e3._menu);
            if (n2.includes(e3._element) || "inside" === e3._config.autoClose && !s2 || "outside" === e3._config.autoClose && s2)
              continue;
            if (e3._menu.contains(t2.target) && ("keyup" === t2.type && "Tab" === t2.key || /input|select|option|textarea|form/i.test(t2.target.tagName)))
              continue;
            const o2 = { relatedTarget: e3._element };
            "click" === t2.type && (o2.clickEvent = t2), e3._completeHide(o2);
          }
        }
        static dataApiKeydownHandler(t2) {
          const e2 = /input|textarea/i.test(t2.target.tagName), i2 = "Escape" === t2.key, n2 = [Ei, Ti].includes(t2.key);
          if (!n2 && !i2)
            return;
          if (e2 && !i2)
            return;
          t2.preventDefault();
          const s2 = this.matches(Ii) ? this : z.prev(this, Ii)[0] || z.next(this, Ii)[0] || z.findOne(Ii, t2.delegateTarget.parentNode), o2 = qi.getOrCreateInstance(s2);
          if (n2)
            return t2.stopPropagation(), o2.show(), void o2._selectMenuItem(t2);
          o2._isShown() && (t2.stopPropagation(), o2.hide(), s2.focus());
        }
      }
      N.on(document, Si, Ii, qi.dataApiKeydownHandler), N.on(document, Si, Pi, qi.dataApiKeydownHandler), N.on(document, Li, qi.clearMenus), N.on(document, Di, qi.clearMenus), N.on(document, Li, Ii, function(t2) {
        t2.preventDefault(), qi.getOrCreateInstance(this).toggle();
      }), m(qi);
      const Vi = "backdrop", Ki = "show", Qi = `mousedown.bs.${Vi}`, Xi = { className: "modal-backdrop", clickCallback: null, isAnimated: false, isVisible: true, rootElement: "body" }, Yi = { className: "string", clickCallback: "(function|null)", isAnimated: "boolean", isVisible: "boolean", rootElement: "(element|string)" };
      class Ui extends H {
        constructor(t2) {
          super(), this._config = this._getConfig(t2), this._isAppended = false, this._element = null;
        }
        static get Default() {
          return Xi;
        }
        static get DefaultType() {
          return Yi;
        }
        static get NAME() {
          return Vi;
        }
        show(t2) {
          if (!this._config.isVisible)
            return void g(t2);
          this._append();
          const e2 = this._getElement();
          this._config.isAnimated && d(e2), e2.classList.add(Ki), this._emulateAnimation(() => {
            g(t2);
          });
        }
        hide(t2) {
          this._config.isVisible ? (this._getElement().classList.remove(Ki), this._emulateAnimation(() => {
            this.dispose(), g(t2);
          })) : g(t2);
        }
        dispose() {
          this._isAppended && (N.off(this._element, Qi), this._element.remove(), this._isAppended = false);
        }
        _getElement() {
          if (!this._element) {
            const t2 = document.createElement("div");
            t2.className = this._config.className, this._config.isAnimated && t2.classList.add("fade"), this._element = t2;
          }
          return this._element;
        }
        _configAfterMerge(t2) {
          return t2.rootElement = r(t2.rootElement), t2;
        }
        _append() {
          if (this._isAppended)
            return;
          const t2 = this._getElement();
          this._config.rootElement.append(t2), N.on(t2, Qi, () => {
            g(this._config.clickCallback);
          }), this._isAppended = true;
        }
        _emulateAnimation(t2) {
          _(t2, this._getElement(), this._config.isAnimated);
        }
      }
      const Gi = ".bs.focustrap", Ji = `focusin${Gi}`, Zi = `keydown.tab${Gi}`, tn = "backward", en = { autofocus: true, trapElement: null }, nn = { autofocus: "boolean", trapElement: "element" };
      class sn extends H {
        constructor(t2) {
          super(), this._config = this._getConfig(t2), this._isActive = false, this._lastTabNavDirection = null;
        }
        static get Default() {
          return en;
        }
        static get DefaultType() {
          return nn;
        }
        static get NAME() {
          return "focustrap";
        }
        activate() {
          this._isActive || (this._config.autofocus && this._config.trapElement.focus(), N.off(document, Gi), N.on(document, Ji, (t2) => this._handleFocusin(t2)), N.on(document, Zi, (t2) => this._handleKeydown(t2)), this._isActive = true);
        }
        deactivate() {
          this._isActive && (this._isActive = false, N.off(document, Gi));
        }
        _handleFocusin(t2) {
          const { trapElement: e2 } = this._config;
          if (t2.target === document || t2.target === e2 || e2.contains(t2.target))
            return;
          const i2 = z.focusableChildren(e2);
          0 === i2.length ? e2.focus() : this._lastTabNavDirection === tn ? i2[i2.length - 1].focus() : i2[0].focus();
        }
        _handleKeydown(t2) {
          "Tab" === t2.key && (this._lastTabNavDirection = t2.shiftKey ? tn : "forward");
        }
      }
      const on = ".fixed-top, .fixed-bottom, .is-fixed, .sticky-top", rn = ".sticky-top", an = "padding-right", ln = "margin-right";
      class cn {
        constructor() {
          this._element = document.body;
        }
        getWidth() {
          const t2 = document.documentElement.clientWidth;
          return Math.abs(window.innerWidth - t2);
        }
        hide() {
          const t2 = this.getWidth();
          this._disableOverFlow(), this._setElementAttributes(this._element, an, (e2) => e2 + t2), this._setElementAttributes(on, an, (e2) => e2 + t2), this._setElementAttributes(rn, ln, (e2) => e2 - t2);
        }
        reset() {
          this._resetElementAttributes(this._element, "overflow"), this._resetElementAttributes(this._element, an), this._resetElementAttributes(on, an), this._resetElementAttributes(rn, ln);
        }
        isOverflowing() {
          return this.getWidth() > 0;
        }
        _disableOverFlow() {
          this._saveInitialAttribute(this._element, "overflow"), this._element.style.overflow = "hidden";
        }
        _setElementAttributes(t2, e2, i2) {
          const n2 = this.getWidth();
          this._applyManipulationCallback(t2, (t3) => {
            if (t3 !== this._element && window.innerWidth > t3.clientWidth + n2)
              return;
            this._saveInitialAttribute(t3, e2);
            const s2 = window.getComputedStyle(t3).getPropertyValue(e2);
            t3.style.setProperty(e2, `${i2(Number.parseFloat(s2))}px`);
          });
        }
        _saveInitialAttribute(t2, e2) {
          const i2 = t2.style.getPropertyValue(e2);
          i2 && F.setDataAttribute(t2, e2, i2);
        }
        _resetElementAttributes(t2, e2) {
          this._applyManipulationCallback(t2, (t3) => {
            const i2 = F.getDataAttribute(t3, e2);
            null !== i2 ? (F.removeDataAttribute(t3, e2), t3.style.setProperty(e2, i2)) : t3.style.removeProperty(e2);
          });
        }
        _applyManipulationCallback(t2, e2) {
          if (o(t2))
            e2(t2);
          else
            for (const i2 of z.find(t2, this._element))
              e2(i2);
        }
      }
      const hn = ".bs.modal", dn = `hide${hn}`, un = `hidePrevented${hn}`, fn = `hidden${hn}`, pn = `show${hn}`, mn = `shown${hn}`, gn = `resize${hn}`, _n = `click.dismiss${hn}`, bn = `mousedown.dismiss${hn}`, vn = `keydown.dismiss${hn}`, yn = `click${hn}.data-api`, wn = "modal-open", An = "show", En = "modal-static", Tn = { backdrop: true, focus: true, keyboard: true }, Cn = { backdrop: "(boolean|string)", focus: "boolean", keyboard: "boolean" };
      class On extends W {
        constructor(t2, e2) {
          super(t2, e2), this._dialog = z.findOne(".modal-dialog", this._element), this._backdrop = this._initializeBackDrop(), this._focustrap = this._initializeFocusTrap(), this._isShown = false, this._isTransitioning = false, this._scrollBar = new cn(), this._addEventListeners();
        }
        static get Default() {
          return Tn;
        }
        static get DefaultType() {
          return Cn;
        }
        static get NAME() {
          return "modal";
        }
        toggle(t2) {
          return this._isShown ? this.hide() : this.show(t2);
        }
        show(t2) {
          this._isShown || this._isTransitioning || N.trigger(this._element, pn, { relatedTarget: t2 }).defaultPrevented || (this._isShown = true, this._isTransitioning = true, this._scrollBar.hide(), document.body.classList.add(wn), this._adjustDialog(), this._backdrop.show(() => this._showElement(t2)));
        }
        hide() {
          this._isShown && !this._isTransitioning && (N.trigger(this._element, dn).defaultPrevented || (this._isShown = false, this._isTransitioning = true, this._focustrap.deactivate(), this._element.classList.remove(An), this._queueCallback(() => this._hideModal(), this._element, this._isAnimated())));
        }
        dispose() {
          N.off(window, hn), N.off(this._dialog, hn), this._backdrop.dispose(), this._focustrap.deactivate(), super.dispose();
        }
        handleUpdate() {
          this._adjustDialog();
        }
        _initializeBackDrop() {
          return new Ui({ isVisible: Boolean(this._config.backdrop), isAnimated: this._isAnimated() });
        }
        _initializeFocusTrap() {
          return new sn({ trapElement: this._element });
        }
        _showElement(t2) {
          document.body.contains(this._element) || document.body.append(this._element), this._element.style.display = "block", this._element.removeAttribute("aria-hidden"), this._element.setAttribute("aria-modal", true), this._element.setAttribute("role", "dialog"), this._element.scrollTop = 0;
          const e2 = z.findOne(".modal-body", this._dialog);
          e2 && (e2.scrollTop = 0), d(this._element), this._element.classList.add(An), this._queueCallback(() => {
            this._config.focus && this._focustrap.activate(), this._isTransitioning = false, N.trigger(this._element, mn, { relatedTarget: t2 });
          }, this._dialog, this._isAnimated());
        }
        _addEventListeners() {
          N.on(this._element, vn, (t2) => {
            "Escape" === t2.key && (this._config.keyboard ? this.hide() : this._triggerBackdropTransition());
          }), N.on(window, gn, () => {
            this._isShown && !this._isTransitioning && this._adjustDialog();
          }), N.on(this._element, bn, (t2) => {
            N.one(this._element, _n, (e2) => {
              this._element === t2.target && this._element === e2.target && ("static" !== this._config.backdrop ? this._config.backdrop && this.hide() : this._triggerBackdropTransition());
            });
          });
        }
        _hideModal() {
          this._element.style.display = "none", this._element.setAttribute("aria-hidden", true), this._element.removeAttribute("aria-modal"), this._element.removeAttribute("role"), this._isTransitioning = false, this._backdrop.hide(() => {
            document.body.classList.remove(wn), this._resetAdjustments(), this._scrollBar.reset(), N.trigger(this._element, fn);
          });
        }
        _isAnimated() {
          return this._element.classList.contains("fade");
        }
        _triggerBackdropTransition() {
          if (N.trigger(this._element, un).defaultPrevented)
            return;
          const t2 = this._element.scrollHeight > document.documentElement.clientHeight, e2 = this._element.style.overflowY;
          "hidden" === e2 || this._element.classList.contains(En) || (t2 || (this._element.style.overflowY = "hidden"), this._element.classList.add(En), this._queueCallback(() => {
            this._element.classList.remove(En), this._queueCallback(() => {
              this._element.style.overflowY = e2;
            }, this._dialog);
          }, this._dialog), this._element.focus());
        }
        _adjustDialog() {
          const t2 = this._element.scrollHeight > document.documentElement.clientHeight, e2 = this._scrollBar.getWidth(), i2 = e2 > 0;
          if (i2 && !t2) {
            const t3 = p() ? "paddingLeft" : "paddingRight";
            this._element.style[t3] = `${e2}px`;
          }
          if (!i2 && t2) {
            const t3 = p() ? "paddingRight" : "paddingLeft";
            this._element.style[t3] = `${e2}px`;
          }
        }
        _resetAdjustments() {
          this._element.style.paddingLeft = "", this._element.style.paddingRight = "";
        }
        static jQueryInterface(t2, e2) {
          return this.each(function() {
            const i2 = On.getOrCreateInstance(this, t2);
            if ("string" == typeof t2) {
              if (void 0 === i2[t2])
                throw new TypeError(`No method named "${t2}"`);
              i2[t2](e2);
            }
          });
        }
      }
      N.on(document, yn, '[data-bs-toggle="modal"]', function(t2) {
        const e2 = z.getElementFromSelector(this);
        ["A", "AREA"].includes(this.tagName) && t2.preventDefault(), N.one(e2, pn, (t3) => {
          t3.defaultPrevented || N.one(e2, fn, () => {
            a(this) && this.focus();
          });
        });
        const i2 = z.findOne(".modal.show");
        i2 && On.getInstance(i2).hide(), On.getOrCreateInstance(e2).toggle(this);
      }), R(On), m(On);
      const xn = ".bs.offcanvas", kn = ".data-api", Ln = `load${xn}${kn}`, Sn = "show", Dn = "showing", $n = "hiding", In = ".offcanvas.show", Nn = `show${xn}`, Pn = `shown${xn}`, Mn = `hide${xn}`, jn = `hidePrevented${xn}`, Fn = `hidden${xn}`, Hn = `resize${xn}`, Wn = `click${xn}${kn}`, Bn = `keydown.dismiss${xn}`, zn = { backdrop: true, keyboard: true, scroll: false }, Rn = { backdrop: "(boolean|string)", keyboard: "boolean", scroll: "boolean" };
      class qn extends W {
        constructor(t2, e2) {
          super(t2, e2), this._isShown = false, this._backdrop = this._initializeBackDrop(), this._focustrap = this._initializeFocusTrap(), this._addEventListeners();
        }
        static get Default() {
          return zn;
        }
        static get DefaultType() {
          return Rn;
        }
        static get NAME() {
          return "offcanvas";
        }
        toggle(t2) {
          return this._isShown ? this.hide() : this.show(t2);
        }
        show(t2) {
          this._isShown || N.trigger(this._element, Nn, { relatedTarget: t2 }).defaultPrevented || (this._isShown = true, this._backdrop.show(), this._config.scroll || new cn().hide(), this._element.setAttribute("aria-modal", true), this._element.setAttribute("role", "dialog"), this._element.classList.add(Dn), this._queueCallback(() => {
            this._config.scroll && !this._config.backdrop || this._focustrap.activate(), this._element.classList.add(Sn), this._element.classList.remove(Dn), N.trigger(this._element, Pn, { relatedTarget: t2 });
          }, this._element, true));
        }
        hide() {
          this._isShown && (N.trigger(this._element, Mn).defaultPrevented || (this._focustrap.deactivate(), this._element.blur(), this._isShown = false, this._element.classList.add($n), this._backdrop.hide(), this._queueCallback(() => {
            this._element.classList.remove(Sn, $n), this._element.removeAttribute("aria-modal"), this._element.removeAttribute("role"), this._config.scroll || new cn().reset(), N.trigger(this._element, Fn);
          }, this._element, true)));
        }
        dispose() {
          this._backdrop.dispose(), this._focustrap.deactivate(), super.dispose();
        }
        _initializeBackDrop() {
          const t2 = Boolean(this._config.backdrop);
          return new Ui({ className: "offcanvas-backdrop", isVisible: t2, isAnimated: true, rootElement: this._element.parentNode, clickCallback: t2 ? () => {
            "static" !== this._config.backdrop ? this.hide() : N.trigger(this._element, jn);
          } : null });
        }
        _initializeFocusTrap() {
          return new sn({ trapElement: this._element });
        }
        _addEventListeners() {
          N.on(this._element, Bn, (t2) => {
            "Escape" === t2.key && (this._config.keyboard ? this.hide() : N.trigger(this._element, jn));
          });
        }
        static jQueryInterface(t2) {
          return this.each(function() {
            const e2 = qn.getOrCreateInstance(this, t2);
            if ("string" == typeof t2) {
              if (void 0 === e2[t2] || t2.startsWith("_") || "constructor" === t2)
                throw new TypeError(`No method named "${t2}"`);
              e2[t2](this);
            }
          });
        }
      }
      N.on(document, Wn, '[data-bs-toggle="offcanvas"]', function(t2) {
        const e2 = z.getElementFromSelector(this);
        if (["A", "AREA"].includes(this.tagName) && t2.preventDefault(), l(this))
          return;
        N.one(e2, Fn, () => {
          a(this) && this.focus();
        });
        const i2 = z.findOne(In);
        i2 && i2 !== e2 && qn.getInstance(i2).hide(), qn.getOrCreateInstance(e2).toggle(this);
      }), N.on(window, Ln, () => {
        for (const t2 of z.find(In))
          qn.getOrCreateInstance(t2).show();
      }), N.on(window, Hn, () => {
        for (const t2 of z.find("[aria-modal][class*=show][class*=offcanvas-]"))
          "fixed" !== getComputedStyle(t2).position && qn.getOrCreateInstance(t2).hide();
      }), R(qn), m(qn);
      const Vn = { "*": ["class", "dir", "id", "lang", "role", /^aria-[\w-]*$/i], a: ["target", "href", "title", "rel"], area: [], b: [], br: [], col: [], code: [], div: [], em: [], hr: [], h1: [], h2: [], h3: [], h4: [], h5: [], h6: [], i: [], img: ["src", "srcset", "alt", "title", "width", "height"], li: [], ol: [], p: [], pre: [], s: [], small: [], span: [], sub: [], sup: [], strong: [], u: [], ul: [] }, Kn = /* @__PURE__ */ new Set(["background", "cite", "href", "itemtype", "longdesc", "poster", "src", "xlink:href"]), Qn = /^(?!javascript:)(?:[a-z0-9+.-]+:|[^&:/?#]*(?:[/?#]|$))/i, Xn = (t2, e2) => {
        const i2 = t2.nodeName.toLowerCase();
        return e2.includes(i2) ? !Kn.has(i2) || Boolean(Qn.test(t2.nodeValue)) : e2.filter((t3) => t3 instanceof RegExp).some((t3) => t3.test(i2));
      }, Yn = { allowList: Vn, content: {}, extraClass: "", html: false, sanitize: true, sanitizeFn: null, template: "<div></div>" }, Un = { allowList: "object", content: "object", extraClass: "(string|function)", html: "boolean", sanitize: "boolean", sanitizeFn: "(null|function)", template: "string" }, Gn = { entry: "(string|element|function|null)", selector: "(string|element)" };
      class Jn extends H {
        constructor(t2) {
          super(), this._config = this._getConfig(t2);
        }
        static get Default() {
          return Yn;
        }
        static get DefaultType() {
          return Un;
        }
        static get NAME() {
          return "TemplateFactory";
        }
        getContent() {
          return Object.values(this._config.content).map((t2) => this._resolvePossibleFunction(t2)).filter(Boolean);
        }
        hasContent() {
          return this.getContent().length > 0;
        }
        changeContent(t2) {
          return this._checkContent(t2), this._config.content = { ...this._config.content, ...t2 }, this;
        }
        toHtml() {
          const t2 = document.createElement("div");
          t2.innerHTML = this._maybeSanitize(this._config.template);
          for (const [e3, i3] of Object.entries(this._config.content))
            this._setContent(t2, i3, e3);
          const e2 = t2.children[0], i2 = this._resolvePossibleFunction(this._config.extraClass);
          return i2 && e2.classList.add(...i2.split(" ")), e2;
        }
        _typeCheckConfig(t2) {
          super._typeCheckConfig(t2), this._checkContent(t2.content);
        }
        _checkContent(t2) {
          for (const [e2, i2] of Object.entries(t2))
            super._typeCheckConfig({ selector: e2, entry: i2 }, Gn);
        }
        _setContent(t2, e2, i2) {
          const n2 = z.findOne(i2, t2);
          n2 && ((e2 = this._resolvePossibleFunction(e2)) ? o(e2) ? this._putElementInTemplate(r(e2), n2) : this._config.html ? n2.innerHTML = this._maybeSanitize(e2) : n2.textContent = e2 : n2.remove());
        }
        _maybeSanitize(t2) {
          return this._config.sanitize ? function(t3, e2, i2) {
            if (!t3.length)
              return t3;
            if (i2 && "function" == typeof i2)
              return i2(t3);
            const n2 = new window.DOMParser().parseFromString(t3, "text/html"), s2 = [].concat(...n2.body.querySelectorAll("*"));
            for (const t4 of s2) {
              const i3 = t4.nodeName.toLowerCase();
              if (!Object.keys(e2).includes(i3)) {
                t4.remove();
                continue;
              }
              const n3 = [].concat(...t4.attributes), s3 = [].concat(e2["*"] || [], e2[i3] || []);
              for (const e3 of n3)
                Xn(e3, s3) || t4.removeAttribute(e3.nodeName);
            }
            return n2.body.innerHTML;
          }(t2, this._config.allowList, this._config.sanitizeFn) : t2;
        }
        _resolvePossibleFunction(t2) {
          return g(t2, [this]);
        }
        _putElementInTemplate(t2, e2) {
          if (this._config.html)
            return e2.innerHTML = "", void e2.append(t2);
          e2.textContent = t2.textContent;
        }
      }
      const Zn = /* @__PURE__ */ new Set(["sanitize", "allowList", "sanitizeFn"]), ts = "fade", es = "show", is = ".modal", ns = "hide.bs.modal", ss = "hover", os = "focus", rs = { AUTO: "auto", TOP: "top", RIGHT: p() ? "left" : "right", BOTTOM: "bottom", LEFT: p() ? "right" : "left" }, as = { allowList: Vn, animation: true, boundary: "clippingParents", container: false, customClass: "", delay: 0, fallbackPlacements: ["top", "right", "bottom", "left"], html: false, offset: [0, 6], placement: "top", popperConfig: null, sanitize: true, sanitizeFn: null, selector: false, template: '<div class="tooltip" role="tooltip"><div class="tooltip-arrow"></div><div class="tooltip-inner"></div></div>', title: "", trigger: "hover focus" }, ls = { allowList: "object", animation: "boolean", boundary: "(string|element)", container: "(string|element|boolean)", customClass: "(string|function)", delay: "(number|object)", fallbackPlacements: "array", html: "boolean", offset: "(array|string|function)", placement: "(string|function)", popperConfig: "(null|object|function)", sanitize: "boolean", sanitizeFn: "(null|function)", selector: "(string|boolean)", template: "string", title: "(string|element|function)", trigger: "string" };
      class cs extends W {
        constructor(t2, e2) {
          if (void 0 === vi)
            throw new TypeError("Bootstrap's tooltips require Popper (https://popper.js.org)");
          super(t2, e2), this._isEnabled = true, this._timeout = 0, this._isHovered = null, this._activeTrigger = {}, this._popper = null, this._templateFactory = null, this._newContent = null, this.tip = null, this._setListeners(), this._config.selector || this._fixTitle();
        }
        static get Default() {
          return as;
        }
        static get DefaultType() {
          return ls;
        }
        static get NAME() {
          return "tooltip";
        }
        enable() {
          this._isEnabled = true;
        }
        disable() {
          this._isEnabled = false;
        }
        toggleEnabled() {
          this._isEnabled = !this._isEnabled;
        }
        toggle() {
          this._isEnabled && (this._activeTrigger.click = !this._activeTrigger.click, this._isShown() ? this._leave() : this._enter());
        }
        dispose() {
          clearTimeout(this._timeout), N.off(this._element.closest(is), ns, this._hideModalHandler), this._element.getAttribute("data-bs-original-title") && this._element.setAttribute("title", this._element.getAttribute("data-bs-original-title")), this._disposePopper(), super.dispose();
        }
        show() {
          if ("none" === this._element.style.display)
            throw new Error("Please use show on visible elements");
          if (!this._isWithContent() || !this._isEnabled)
            return;
          const t2 = N.trigger(this._element, this.constructor.eventName("show")), e2 = (c(this._element) || this._element.ownerDocument.documentElement).contains(this._element);
          if (t2.defaultPrevented || !e2)
            return;
          this._disposePopper();
          const i2 = this._getTipElement();
          this._element.setAttribute("aria-describedby", i2.getAttribute("id"));
          const { container: n2 } = this._config;
          if (this._element.ownerDocument.documentElement.contains(this.tip) || (n2.append(i2), N.trigger(this._element, this.constructor.eventName("inserted"))), this._popper = this._createPopper(i2), i2.classList.add(es), "ontouchstart" in document.documentElement)
            for (const t3 of [].concat(...document.body.children))
              N.on(t3, "mouseover", h);
          this._queueCallback(() => {
            N.trigger(this._element, this.constructor.eventName("shown")), false === this._isHovered && this._leave(), this._isHovered = false;
          }, this.tip, this._isAnimated());
        }
        hide() {
          if (this._isShown() && !N.trigger(this._element, this.constructor.eventName("hide")).defaultPrevented) {
            if (this._getTipElement().classList.remove(es), "ontouchstart" in document.documentElement)
              for (const t2 of [].concat(...document.body.children))
                N.off(t2, "mouseover", h);
            this._activeTrigger.click = false, this._activeTrigger[os] = false, this._activeTrigger[ss] = false, this._isHovered = null, this._queueCallback(() => {
              this._isWithActiveTrigger() || (this._isHovered || this._disposePopper(), this._element.removeAttribute("aria-describedby"), N.trigger(this._element, this.constructor.eventName("hidden")));
            }, this.tip, this._isAnimated());
          }
        }
        update() {
          this._popper && this._popper.update();
        }
        _isWithContent() {
          return Boolean(this._getTitle());
        }
        _getTipElement() {
          return this.tip || (this.tip = this._createTipElement(this._newContent || this._getContentForTemplate())), this.tip;
        }
        _createTipElement(t2) {
          const e2 = this._getTemplateFactory(t2).toHtml();
          if (!e2)
            return null;
          e2.classList.remove(ts, es), e2.classList.add(`bs-${this.constructor.NAME}-auto`);
          const i2 = ((t3) => {
            do {
              t3 += Math.floor(1e6 * Math.random());
            } while (document.getElementById(t3));
            return t3;
          })(this.constructor.NAME).toString();
          return e2.setAttribute("id", i2), this._isAnimated() && e2.classList.add(ts), e2;
        }
        setContent(t2) {
          this._newContent = t2, this._isShown() && (this._disposePopper(), this.show());
        }
        _getTemplateFactory(t2) {
          return this._templateFactory ? this._templateFactory.changeContent(t2) : this._templateFactory = new Jn({ ...this._config, content: t2, extraClass: this._resolvePossibleFunction(this._config.customClass) }), this._templateFactory;
        }
        _getContentForTemplate() {
          return { ".tooltip-inner": this._getTitle() };
        }
        _getTitle() {
          return this._resolvePossibleFunction(this._config.title) || this._element.getAttribute("data-bs-original-title");
        }
        _initializeOnDelegatedTarget(t2) {
          return this.constructor.getOrCreateInstance(t2.delegateTarget, this._getDelegateConfig());
        }
        _isAnimated() {
          return this._config.animation || this.tip && this.tip.classList.contains(ts);
        }
        _isShown() {
          return this.tip && this.tip.classList.contains(es);
        }
        _createPopper(t2) {
          const e2 = g(this._config.placement, [this, t2, this._element]), i2 = rs[e2.toUpperCase()];
          return bi(this._element, t2, this._getPopperConfig(i2));
        }
        _getOffset() {
          const { offset: t2 } = this._config;
          return "string" == typeof t2 ? t2.split(",").map((t3) => Number.parseInt(t3, 10)) : "function" == typeof t2 ? (e2) => t2(e2, this._element) : t2;
        }
        _resolvePossibleFunction(t2) {
          return g(t2, [this._element]);
        }
        _getPopperConfig(t2) {
          const e2 = { placement: t2, modifiers: [{ name: "flip", options: { fallbackPlacements: this._config.fallbackPlacements } }, { name: "offset", options: { offset: this._getOffset() } }, { name: "preventOverflow", options: { boundary: this._config.boundary } }, { name: "arrow", options: { element: `.${this.constructor.NAME}-arrow` } }, { name: "preSetPlacement", enabled: true, phase: "beforeMain", fn: (t3) => {
            this._getTipElement().setAttribute("data-popper-placement", t3.state.placement);
          } }] };
          return { ...e2, ...g(this._config.popperConfig, [e2]) };
        }
        _setListeners() {
          const t2 = this._config.trigger.split(" ");
          for (const e2 of t2)
            if ("click" === e2)
              N.on(this._element, this.constructor.eventName("click"), this._config.selector, (t3) => {
                this._initializeOnDelegatedTarget(t3).toggle();
              });
            else if ("manual" !== e2) {
              const t3 = e2 === ss ? this.constructor.eventName("mouseenter") : this.constructor.eventName("focusin"), i2 = e2 === ss ? this.constructor.eventName("mouseleave") : this.constructor.eventName("focusout");
              N.on(this._element, t3, this._config.selector, (t4) => {
                const e3 = this._initializeOnDelegatedTarget(t4);
                e3._activeTrigger["focusin" === t4.type ? os : ss] = true, e3._enter();
              }), N.on(this._element, i2, this._config.selector, (t4) => {
                const e3 = this._initializeOnDelegatedTarget(t4);
                e3._activeTrigger["focusout" === t4.type ? os : ss] = e3._element.contains(t4.relatedTarget), e3._leave();
              });
            }
          this._hideModalHandler = () => {
            this._element && this.hide();
          }, N.on(this._element.closest(is), ns, this._hideModalHandler);
        }
        _fixTitle() {
          const t2 = this._element.getAttribute("title");
          t2 && (this._element.getAttribute("aria-label") || this._element.textContent.trim() || this._element.setAttribute("aria-label", t2), this._element.setAttribute("data-bs-original-title", t2), this._element.removeAttribute("title"));
        }
        _enter() {
          this._isShown() || this._isHovered ? this._isHovered = true : (this._isHovered = true, this._setTimeout(() => {
            this._isHovered && this.show();
          }, this._config.delay.show));
        }
        _leave() {
          this._isWithActiveTrigger() || (this._isHovered = false, this._setTimeout(() => {
            this._isHovered || this.hide();
          }, this._config.delay.hide));
        }
        _setTimeout(t2, e2) {
          clearTimeout(this._timeout), this._timeout = setTimeout(t2, e2);
        }
        _isWithActiveTrigger() {
          return Object.values(this._activeTrigger).includes(true);
        }
        _getConfig(t2) {
          const e2 = F.getDataAttributes(this._element);
          for (const t3 of Object.keys(e2))
            Zn.has(t3) && delete e2[t3];
          return t2 = { ...e2, ..."object" == typeof t2 && t2 ? t2 : {} }, t2 = this._mergeConfigObj(t2), t2 = this._configAfterMerge(t2), this._typeCheckConfig(t2), t2;
        }
        _configAfterMerge(t2) {
          return t2.container = false === t2.container ? document.body : r(t2.container), "number" == typeof t2.delay && (t2.delay = { show: t2.delay, hide: t2.delay }), "number" == typeof t2.title && (t2.title = t2.title.toString()), "number" == typeof t2.content && (t2.content = t2.content.toString()), t2;
        }
        _getDelegateConfig() {
          const t2 = {};
          for (const [e2, i2] of Object.entries(this._config))
            this.constructor.Default[e2] !== i2 && (t2[e2] = i2);
          return t2.selector = false, t2.trigger = "manual", t2;
        }
        _disposePopper() {
          this._popper && (this._popper.destroy(), this._popper = null), this.tip && (this.tip.remove(), this.tip = null);
        }
        static jQueryInterface(t2) {
          return this.each(function() {
            const e2 = cs.getOrCreateInstance(this, t2);
            if ("string" == typeof t2) {
              if (void 0 === e2[t2])
                throw new TypeError(`No method named "${t2}"`);
              e2[t2]();
            }
          });
        }
      }
      m(cs);
      const hs = { ...cs.Default, content: "", offset: [0, 8], placement: "right", template: '<div class="popover" role="tooltip"><div class="popover-arrow"></div><h3 class="popover-header"></h3><div class="popover-body"></div></div>', trigger: "click" }, ds = { ...cs.DefaultType, content: "(null|string|element|function)" };
      class us extends cs {
        static get Default() {
          return hs;
        }
        static get DefaultType() {
          return ds;
        }
        static get NAME() {
          return "popover";
        }
        _isWithContent() {
          return this._getTitle() || this._getContent();
        }
        _getContentForTemplate() {
          return { ".popover-header": this._getTitle(), ".popover-body": this._getContent() };
        }
        _getContent() {
          return this._resolvePossibleFunction(this._config.content);
        }
        static jQueryInterface(t2) {
          return this.each(function() {
            const e2 = us.getOrCreateInstance(this, t2);
            if ("string" == typeof t2) {
              if (void 0 === e2[t2])
                throw new TypeError(`No method named "${t2}"`);
              e2[t2]();
            }
          });
        }
      }
      m(us);
      const fs = ".bs.scrollspy", ps = `activate${fs}`, ms = `click${fs}`, gs = `load${fs}.data-api`, _s = "active", bs = "[href]", vs = ".nav-link", ys = `${vs}, .nav-item > ${vs}, .list-group-item`, ws = { offset: null, rootMargin: "0px 0px -25%", smoothScroll: false, target: null, threshold: [0.1, 0.5, 1] }, As = { offset: "(number|null)", rootMargin: "string", smoothScroll: "boolean", target: "element", threshold: "array" };
      class Es extends W {
        constructor(t2, e2) {
          super(t2, e2), this._targetLinks = /* @__PURE__ */ new Map(), this._observableSections = /* @__PURE__ */ new Map(), this._rootElement = "visible" === getComputedStyle(this._element).overflowY ? null : this._element, this._activeTarget = null, this._observer = null, this._previousScrollData = { visibleEntryTop: 0, parentScrollTop: 0 }, this.refresh();
        }
        static get Default() {
          return ws;
        }
        static get DefaultType() {
          return As;
        }
        static get NAME() {
          return "scrollspy";
        }
        refresh() {
          this._initializeTargetsAndObservables(), this._maybeEnableSmoothScroll(), this._observer ? this._observer.disconnect() : this._observer = this._getNewObserver();
          for (const t2 of this._observableSections.values())
            this._observer.observe(t2);
        }
        dispose() {
          this._observer.disconnect(), super.dispose();
        }
        _configAfterMerge(t2) {
          return t2.target = r(t2.target) || document.body, t2.rootMargin = t2.offset ? `${t2.offset}px 0px -30%` : t2.rootMargin, "string" == typeof t2.threshold && (t2.threshold = t2.threshold.split(",").map((t3) => Number.parseFloat(t3))), t2;
        }
        _maybeEnableSmoothScroll() {
          this._config.smoothScroll && (N.off(this._config.target, ms), N.on(this._config.target, ms, bs, (t2) => {
            const e2 = this._observableSections.get(t2.target.hash);
            if (e2) {
              t2.preventDefault();
              const i2 = this._rootElement || window, n2 = e2.offsetTop - this._element.offsetTop;
              if (i2.scrollTo)
                return void i2.scrollTo({ top: n2, behavior: "smooth" });
              i2.scrollTop = n2;
            }
          }));
        }
        _getNewObserver() {
          const t2 = { root: this._rootElement, threshold: this._config.threshold, rootMargin: this._config.rootMargin };
          return new IntersectionObserver((t3) => this._observerCallback(t3), t2);
        }
        _observerCallback(t2) {
          const e2 = (t3) => this._targetLinks.get(`#${t3.target.id}`), i2 = (t3) => {
            this._previousScrollData.visibleEntryTop = t3.target.offsetTop, this._process(e2(t3));
          }, n2 = (this._rootElement || document.documentElement).scrollTop, s2 = n2 >= this._previousScrollData.parentScrollTop;
          this._previousScrollData.parentScrollTop = n2;
          for (const o2 of t2) {
            if (!o2.isIntersecting) {
              this._activeTarget = null, this._clearActiveClass(e2(o2));
              continue;
            }
            const t3 = o2.target.offsetTop >= this._previousScrollData.visibleEntryTop;
            if (s2 && t3) {
              if (i2(o2), !n2)
                return;
            } else
              s2 || t3 || i2(o2);
          }
        }
        _initializeTargetsAndObservables() {
          this._targetLinks = /* @__PURE__ */ new Map(), this._observableSections = /* @__PURE__ */ new Map();
          const t2 = z.find(bs, this._config.target);
          for (const e2 of t2) {
            if (!e2.hash || l(e2))
              continue;
            const t3 = z.findOne(decodeURI(e2.hash), this._element);
            a(t3) && (this._targetLinks.set(decodeURI(e2.hash), e2), this._observableSections.set(e2.hash, t3));
          }
        }
        _process(t2) {
          this._activeTarget !== t2 && (this._clearActiveClass(this._config.target), this._activeTarget = t2, t2.classList.add(_s), this._activateParents(t2), N.trigger(this._element, ps, { relatedTarget: t2 }));
        }
        _activateParents(t2) {
          if (t2.classList.contains("dropdown-item"))
            z.findOne(".dropdown-toggle", t2.closest(".dropdown")).classList.add(_s);
          else
            for (const e2 of z.parents(t2, ".nav, .list-group"))
              for (const t3 of z.prev(e2, ys))
                t3.classList.add(_s);
        }
        _clearActiveClass(t2) {
          t2.classList.remove(_s);
          const e2 = z.find(`${bs}.${_s}`, t2);
          for (const t3 of e2)
            t3.classList.remove(_s);
        }
        static jQueryInterface(t2) {
          return this.each(function() {
            const e2 = Es.getOrCreateInstance(this, t2);
            if ("string" == typeof t2) {
              if (void 0 === e2[t2] || t2.startsWith("_") || "constructor" === t2)
                throw new TypeError(`No method named "${t2}"`);
              e2[t2]();
            }
          });
        }
      }
      N.on(window, gs, () => {
        for (const t2 of z.find('[data-bs-spy="scroll"]'))
          Es.getOrCreateInstance(t2);
      }), m(Es);
      const Ts = ".bs.tab", Cs = `hide${Ts}`, Os = `hidden${Ts}`, xs = `show${Ts}`, ks = `shown${Ts}`, Ls = `click${Ts}`, Ss = `keydown${Ts}`, Ds = `load${Ts}`, $s = "ArrowLeft", Is = "ArrowRight", Ns = "ArrowUp", Ps = "ArrowDown", Ms = "Home", js = "End", Fs = "active", Hs = "fade", Ws = "show", Bs = ".dropdown-toggle", zs = `:not(${Bs})`, Rs = '[data-bs-toggle="tab"], [data-bs-toggle="pill"], [data-bs-toggle="list"]', qs = `.nav-link${zs}, .list-group-item${zs}, [role="tab"]${zs}, ${Rs}`, Vs = `.${Fs}[data-bs-toggle="tab"], .${Fs}[data-bs-toggle="pill"], .${Fs}[data-bs-toggle="list"]`;
      class Ks extends W {
        constructor(t2) {
          super(t2), this._parent = this._element.closest('.list-group, .nav, [role="tablist"]'), this._parent && (this._setInitialAttributes(this._parent, this._getChildren()), N.on(this._element, Ss, (t3) => this._keydown(t3)));
        }
        static get NAME() {
          return "tab";
        }
        show() {
          const t2 = this._element;
          if (this._elemIsActive(t2))
            return;
          const e2 = this._getActiveElem(), i2 = e2 ? N.trigger(e2, Cs, { relatedTarget: t2 }) : null;
          N.trigger(t2, xs, { relatedTarget: e2 }).defaultPrevented || i2 && i2.defaultPrevented || (this._deactivate(e2, t2), this._activate(t2, e2));
        }
        _activate(t2, e2) {
          t2 && (t2.classList.add(Fs), this._activate(z.getElementFromSelector(t2)), this._queueCallback(() => {
            "tab" === t2.getAttribute("role") ? (t2.removeAttribute("tabindex"), t2.setAttribute("aria-selected", true), this._toggleDropDown(t2, true), N.trigger(t2, ks, { relatedTarget: e2 })) : t2.classList.add(Ws);
          }, t2, t2.classList.contains(Hs)));
        }
        _deactivate(t2, e2) {
          t2 && (t2.classList.remove(Fs), t2.blur(), this._deactivate(z.getElementFromSelector(t2)), this._queueCallback(() => {
            "tab" === t2.getAttribute("role") ? (t2.setAttribute("aria-selected", false), t2.setAttribute("tabindex", "-1"), this._toggleDropDown(t2, false), N.trigger(t2, Os, { relatedTarget: e2 })) : t2.classList.remove(Ws);
          }, t2, t2.classList.contains(Hs)));
        }
        _keydown(t2) {
          if (![$s, Is, Ns, Ps, Ms, js].includes(t2.key))
            return;
          t2.stopPropagation(), t2.preventDefault();
          const e2 = this._getChildren().filter((t3) => !l(t3));
          let i2;
          if ([Ms, js].includes(t2.key))
            i2 = e2[t2.key === Ms ? 0 : e2.length - 1];
          else {
            const n2 = [Is, Ps].includes(t2.key);
            i2 = b(e2, t2.target, n2, true);
          }
          i2 && (i2.focus({ preventScroll: true }), Ks.getOrCreateInstance(i2).show());
        }
        _getChildren() {
          return z.find(qs, this._parent);
        }
        _getActiveElem() {
          return this._getChildren().find((t2) => this._elemIsActive(t2)) || null;
        }
        _setInitialAttributes(t2, e2) {
          this._setAttributeIfNotExists(t2, "role", "tablist");
          for (const t3 of e2)
            this._setInitialAttributesOnChild(t3);
        }
        _setInitialAttributesOnChild(t2) {
          t2 = this._getInnerElement(t2);
          const e2 = this._elemIsActive(t2), i2 = this._getOuterElement(t2);
          t2.setAttribute("aria-selected", e2), i2 !== t2 && this._setAttributeIfNotExists(i2, "role", "presentation"), e2 || t2.setAttribute("tabindex", "-1"), this._setAttributeIfNotExists(t2, "role", "tab"), this._setInitialAttributesOnTargetPanel(t2);
        }
        _setInitialAttributesOnTargetPanel(t2) {
          const e2 = z.getElementFromSelector(t2);
          e2 && (this._setAttributeIfNotExists(e2, "role", "tabpanel"), t2.id && this._setAttributeIfNotExists(e2, "aria-labelledby", `${t2.id}`));
        }
        _toggleDropDown(t2, e2) {
          const i2 = this._getOuterElement(t2);
          if (!i2.classList.contains("dropdown"))
            return;
          const n2 = (t3, n3) => {
            const s2 = z.findOne(t3, i2);
            s2 && s2.classList.toggle(n3, e2);
          };
          n2(Bs, Fs), n2(".dropdown-menu", Ws), i2.setAttribute("aria-expanded", e2);
        }
        _setAttributeIfNotExists(t2, e2, i2) {
          t2.hasAttribute(e2) || t2.setAttribute(e2, i2);
        }
        _elemIsActive(t2) {
          return t2.classList.contains(Fs);
        }
        _getInnerElement(t2) {
          return t2.matches(qs) ? t2 : z.findOne(qs, t2);
        }
        _getOuterElement(t2) {
          return t2.closest(".nav-item, .list-group-item") || t2;
        }
        static jQueryInterface(t2) {
          return this.each(function() {
            const e2 = Ks.getOrCreateInstance(this);
            if ("string" == typeof t2) {
              if (void 0 === e2[t2] || t2.startsWith("_") || "constructor" === t2)
                throw new TypeError(`No method named "${t2}"`);
              e2[t2]();
            }
          });
        }
      }
      N.on(document, Ls, Rs, function(t2) {
        ["A", "AREA"].includes(this.tagName) && t2.preventDefault(), l(this) || Ks.getOrCreateInstance(this).show();
      }), N.on(window, Ds, () => {
        for (const t2 of z.find(Vs))
          Ks.getOrCreateInstance(t2);
      }), m(Ks);
      const Qs = ".bs.toast", Xs = `mouseover${Qs}`, Ys = `mouseout${Qs}`, Us = `focusin${Qs}`, Gs = `focusout${Qs}`, Js = `hide${Qs}`, Zs = `hidden${Qs}`, to = `show${Qs}`, eo = `shown${Qs}`, io = "hide", no = "show", so = "showing", oo = { animation: "boolean", autohide: "boolean", delay: "number" }, ro = { animation: true, autohide: true, delay: 5e3 };
      class ao extends W {
        constructor(t2, e2) {
          super(t2, e2), this._timeout = null, this._hasMouseInteraction = false, this._hasKeyboardInteraction = false, this._setListeners();
        }
        static get Default() {
          return ro;
        }
        static get DefaultType() {
          return oo;
        }
        static get NAME() {
          return "toast";
        }
        show() {
          N.trigger(this._element, to).defaultPrevented || (this._clearTimeout(), this._config.animation && this._element.classList.add("fade"), this._element.classList.remove(io), d(this._element), this._element.classList.add(no, so), this._queueCallback(() => {
            this._element.classList.remove(so), N.trigger(this._element, eo), this._maybeScheduleHide();
          }, this._element, this._config.animation));
        }
        hide() {
          this.isShown() && (N.trigger(this._element, Js).defaultPrevented || (this._element.classList.add(so), this._queueCallback(() => {
            this._element.classList.add(io), this._element.classList.remove(so, no), N.trigger(this._element, Zs);
          }, this._element, this._config.animation)));
        }
        dispose() {
          this._clearTimeout(), this.isShown() && this._element.classList.remove(no), super.dispose();
        }
        isShown() {
          return this._element.classList.contains(no);
        }
        _maybeScheduleHide() {
          this._config.autohide && (this._hasMouseInteraction || this._hasKeyboardInteraction || (this._timeout = setTimeout(() => {
            this.hide();
          }, this._config.delay)));
        }
        _onInteraction(t2, e2) {
          switch (t2.type) {
            case "mouseover":
            case "mouseout":
              this._hasMouseInteraction = e2;
              break;
            case "focusin":
            case "focusout":
              this._hasKeyboardInteraction = e2;
          }
          if (e2)
            return void this._clearTimeout();
          const i2 = t2.relatedTarget;
          this._element === i2 || this._element.contains(i2) || this._maybeScheduleHide();
        }
        _setListeners() {
          N.on(this._element, Xs, (t2) => this._onInteraction(t2, true)), N.on(this._element, Ys, (t2) => this._onInteraction(t2, false)), N.on(this._element, Us, (t2) => this._onInteraction(t2, true)), N.on(this._element, Gs, (t2) => this._onInteraction(t2, false));
        }
        _clearTimeout() {
          clearTimeout(this._timeout), this._timeout = null;
        }
        static jQueryInterface(t2) {
          return this.each(function() {
            const e2 = ao.getOrCreateInstance(this, t2);
            if ("string" == typeof t2) {
              if (void 0 === e2[t2])
                throw new TypeError(`No method named "${t2}"`);
              e2[t2](this);
            }
          });
        }
      }
      return R(ao), m(ao), { Alert: Q, Button: Y, Carousel: xt, Collapse: Bt, Dropdown: qi, Modal: On, Offcanvas: qn, Popover: us, ScrollSpy: Es, Tab: Ks, Toast: ao, Tooltip: cs };
    });
  }
});
export default require_bootstrap_bundle_min();
/*! Bundled license information:

bootstrap/dist/js/bootstrap.bundle.min.js:
  (*!
    * Bootstrap v5.3.2 (https://getbootstrap.com/)
    * Copyright 2011-2023 The Bootstrap Authors (https://github.com/twbs/bootstrap/graphs/contributors)
    * Licensed under MIT (https://github.com/twbs/bootstrap/blob/main/LICENSE)
    *)
*/
//# sourceMappingURL=bootstrap_dist_js_bootstrap__bundle__min.js.map
