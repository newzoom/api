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
      sourceBuffer.appendBuffer(queue.shift());
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
    sourceBuffer.onupdatestart = (e) => {
      console.log("update start", sourceBuffer.buffered.length);
    };
    sourceBuffer.onupdateend = (e) => {
      console.log("update end", sourceBuffer.buffered.length);
    };
    sourceBuffer.onupdate = (e) => {
      console.log("updated", sourceBuffer.buffered.length);
      if (st == undefined) {
        st = Date.now();
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

  if (window["WebSocket"]) {
    ws = new WebSocket("ws://localhost:8080/ws");
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
