package sdk

import (
	"fmt"
	"strings"
)

func GenerateSDK(publicKeyPEM, kid, rpcEndpoint string) string {
	escapedPEM := strings.ReplaceAll(publicKeyPEM, "\n", "\\n")
	return fmt.Sprintf(`(function() {
  "use strict";
  var CONFIG = {
    kid: "%s",
    publicKeyPEM: "%s",
    rpcEndpoint: "%s/rpc"
  };

  function getVisitorID() {
    var id = localStorage.getItem("_tk_vid");
    if (!id) {
      id = crypto.randomUUID();
      localStorage.setItem("_tk_vid", id);
    }
    return id;
  }

  function getSessionID() {
    var id = sessionStorage.getItem("_tk_sid");
    if (!id) {
      id = crypto.randomUUID();
      sessionStorage.setItem("_tk_sid", id);
    }
    return id;
  }

  function base64Encode(buf) {
    var bytes = new Uint8Array(buf);
    var binary = "";
    for (var i = 0; i < bytes.byteLength; i++) {
      binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
  }

  async function importPublicKey(pem) {
    var pemContents = pem.replace("-----BEGIN PUBLIC KEY-----", "")
      .replace("-----END PUBLIC KEY-----", "")
      .replace(/\\n/g, "");
    var binaryDer = Uint8Array.from(atob(pemContents), function(c) { return c.charCodeAt(0); });
    return crypto.subtle.importKey("spki", binaryDer.buffer,
      { name: "RSA-OAEP", hash: "SHA-256" }, false, ["encrypt"]);
  }

  async function encrypt(payload) {
    var pubKey = await importPublicKey(CONFIG.publicKeyPEM);
    var dataKey = crypto.getRandomValues(new Uint8Array(32));
    var ek = await crypto.subtle.encrypt({ name: "RSA-OAEP" }, pubKey, dataKey);
    var aesKey = await crypto.subtle.importKey("raw", dataKey, "AES-GCM", false, ["encrypt"]);
    var nonce = crypto.getRandomValues(new Uint8Array(12));
    var plaintext = new TextEncoder().encode(JSON.stringify(payload));
    var ct = await crypto.subtle.encrypt({ name: "AES-GCM", iv: nonce }, aesKey, plaintext);
    return {
      ek: base64Encode(ek),
      nonce: base64Encode(nonce),
      ct: base64Encode(ct),
      ts: Math.floor(Date.now() / 1000),
      nonce2: crypto.randomUUID(),
      kid: CONFIG.kid
    };
  }

  async function sendRPC(method, params) {
    var encrypted = await encrypt(params);
    var body = JSON.stringify({
      jsonrpc: "2.0",
      method: method,
      params: encrypted,
      id: crypto.randomUUID()
    });
    var resp = await fetch(CONFIG.rpcEndpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: body
    });
    return resp.json();
  }

  var queue = [];
  var flushTimer = null;

  function flushEvents(siteKey) {
    if (queue.length === 0) return;
    var events = queue.splice(0, queue.length);
    sendRPC("track.collectEvents", {
      site_key: siteKey,
      visitor_id: getVisitorID(),
      session_id: getSessionID(),
      events: events
    });
  }

  window.TrackSDK = {
    init: function(siteKey) {
      this._siteKey = siteKey;
      this.trackPageview();
    },
    trackPageview: function() {
      queue.push({
        type: "pageview",
        url: location.href,
        title: document.title,
        referrer: document.referrer
      });
      this._scheduleFlush();
    },
    trackEvent: function(eventType, props) {
      queue.push({
        type: eventType,
        url: location.href,
        title: document.title,
        referrer: document.referrer,
        props: props || {}
      });
      this._scheduleFlush();
    },
    _scheduleFlush: function() {
      var self = this;
      if (flushTimer) clearTimeout(flushTimer);
      flushTimer = setTimeout(function() { flushEvents(self._siteKey); }, 1000);
    }
  };
})();
`, kid, escapedPEM, rpcEndpoint)
}

func GenerateClickPage(token, publicKeyPEM, kid, rpcEndpoint, targetURL string) string {
	escapedPEM := strings.ReplaceAll(publicKeyPEM, "\n", "\\n")
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Redirecting...</title></head>
<body>
<p>Redirecting, please wait...</p>
<script>
(async function() {
  var CONFIG = {
    kid: "%s",
    publicKeyPEM: "%s",
    rpcEndpoint: "%s/rpc",
    token: "%s",
    targetURL: "%s"
  };

  function base64Encode(buf) {
    var bytes = new Uint8Array(buf);
    var binary = "";
    for (var i = 0; i < bytes.byteLength; i++) {
      binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
  }

  async function importPublicKey(pem) {
    var pemContents = pem.replace("-----BEGIN PUBLIC KEY-----", "")
      .replace("-----END PUBLIC KEY-----", "")
      .replace(/\\n/g, "");
    var binaryDer = Uint8Array.from(atob(pemContents), function(c) { return c.charCodeAt(0); });
    return crypto.subtle.importKey("spki", binaryDer.buffer,
      { name: "RSA-OAEP", hash: "SHA-256" }, false, ["encrypt"]);
  }

  try {
    var visitorID = localStorage.getItem("_tk_vid");
    if (!visitorID) {
      visitorID = crypto.randomUUID();
      localStorage.setItem("_tk_vid", visitorID);
    }

    var payload = {
      token: CONFIG.token,
      visitor_id: visitorID,
      env: {
        screen_width: screen.width,
        screen_height: screen.height,
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
        language: navigator.language,
        platform: navigator.platform
      }
    };

    var pubKey = await importPublicKey(CONFIG.publicKeyPEM);
    var dataKey = crypto.getRandomValues(new Uint8Array(32));
    var ek = await crypto.subtle.encrypt({ name: "RSA-OAEP" }, pubKey, dataKey);
    var aesKey = await crypto.subtle.importKey("raw", dataKey, "AES-GCM", false, ["encrypt"]);
    var nonce = crypto.getRandomValues(new Uint8Array(12));
    var plaintext = new TextEncoder().encode(JSON.stringify(payload));
    var ct = await crypto.subtle.encrypt({ name: "AES-GCM", iv: nonce }, aesKey, plaintext);

    var body = JSON.stringify({
      jsonrpc: "2.0",
      method: "track.collectClick",
      params: {
        ek: base64Encode(ek),
        nonce: base64Encode(nonce),
        ct: base64Encode(ct),
        ts: Math.floor(Date.now() / 1000),
        nonce2: crypto.randomUUID(),
        kid: CONFIG.kid
      },
      id: "1"
    });

    await fetch(CONFIG.rpcEndpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: body
    });
  } catch(e) {
    console.error("tracking error:", e);
  }
  window.location.href = CONFIG.targetURL;
})();
</script>
</body>
</html>`, kid, escapedPEM, rpcEndpoint, token, targetURL)
}
