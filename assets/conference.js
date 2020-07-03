window.onload = function () {
  let ws;
  let recorder;
  let mimeType = 'video/webm; codecs="vp8"';
  if (
    !MediaSource.isTypeSupported(mimeType) ||
    !MediaRecorder.isTypeSupported(mimeType)
  ) {
    alert(
      "your browser doesn't support our video quality please come back later for more updates"
    );
    return;
  }
  let queue = [];
  let sourceBuffer;
  let buffer;
  let b;
  let ranges;
  const inputVid = document.querySelector("#input-vid");
  const outputVid = document.querySelector("#output-vid");
  const cvb = document.querySelector("button.capture-button");

  let mediaSource = new MediaSource();
  function appendItemToBuffer() {
    if (
      mediaSource.readyState === "open" &&
      sourceBuffer &&
      !sourceBuffer.updating &&
      queue.length > 0
    ) {
      buffer = queue.shift();
      b = new Int8Array(buffer.slice(0, 37));
      var res = "";
      var i;
      for (i = 1; i < b.length; i++) {
        res += String.fromCharCode(b[i]);
      }
      console.log(b[0], res);
      if (b[0] == 0) {
        sourceBuffer.appendBuffer(buffer.slice(37, buffer.byteLength));
      }
    }

    // Limit the total buffer size to 20 minutes
    // This way we don't run out of RAM
    if (
      outputVid.buffered.length &&
      outputVid.buffered.end(0) - outputVid.buffered.start(0) > 1200
    ) {
      sourceBuffer.remove(0, outputVid.buffered.end(0) - 1200);
    }
  }

  outputVid.src = URL.createObjectURL(mediaSource);
  outputVid.onpaused = (e) => {
    console.log("pause", outputVid.readyState);
  };

  mediaSource.onsourceopen = (e) => {
    sourceBuffer = mediaSource.addSourceBuffer(mimeType);
    sourceBuffer.onerror = (e) => {
      console.log("buffer error");
    };
    sourceBuffer.onupdate = (e) => {
      console.log("updated", sourceBuffer.buffered.length);
      if (outputVid.paused) {
        outputVid.play();
      }
      appendItemToBuffer();
    };
    sourceBuffer.mode = "sequence";
  };

  cvb.onclick = () => {
    navigator.mediaDevices
      .getUserMedia({ video: true })
      .then((stream) => {
        inputVid.srcObject = stream;
        stream.oninactive = (e) => {
          recorder.stop();
        };
        window.stream = stream;
        recorder = new MediaRecorder(stream, { mimeType });
        recorder.onstop = (e) => {
          if (window.stream.active) {
            recorder.start();
          }
        };
        recorder.ondataavailable = (e) => {
          if (e.data && e.data.size > 0) {
            e.data.arrayBuffer().then((b) => {
              ws.send(b);
            });
          }
        };
        recorder.onstart = (e) => {
          setTimeout(() => {
            recorder.stop();
          }, 100);
        };
        recorder.start();
      })
      .catch((err) => {
        console.error(err);
      });
  };

  let sliced;
  let conID = "04a8506e-4e9d-47e5-9410-223df26ac689";
  if (window["WebSocket"]) {
    var uri = `ws://${document.location.host}/ws/${conID}?token=${u.access_token}`;
    ws = new WebSocket(encodeURI(uri));
    ws.binaryType = "arraybuffer";

    ws.onclose = (e) => {
      console.log("connection lost");
    };

    // let newBlob;
    ws.onmessage = (e) => {
      queue.push(e.data);
      appendItemToBuffer();
    };

    ws.onerror = (err) => {
      console.error(err);
    };
  } else {
    alert("Your browser does not support WebSockets.");
  }
};
