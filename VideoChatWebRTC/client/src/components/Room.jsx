import { useEffect, useRef } from "react";
import { useParams } from "react-router-dom";

const Room = () => {
  var params = useParams();

  const userVideo = useRef();
  const userStream = useRef();
  const partnerVideo = useRef();
  const peerRef = useRef();
  const webSocketRef = useRef();

  const openCamera = async () => {
    const allDevices = await navigator.mediaDevices.enumerateDevices();

    const cameras = allDevices.filter((device) => device.kind == "videoinput");
    console.log(cameras);

    const constraint = {
      audio: true,
      video: {
        deviceId: cameras[0].deviceId,
      },
    };

    try {
      return await navigator.mediaDevices.getUserMedia(constraint);
    } catch (err) {
      console.log(err);
    }
  };

  // useEffect(() => {
  //   openCamera().then((stream) => {
  //     userVideo.current.srcObject = stream;
  //     userStream.current = stream;

  //     webSocketRef.current = new WebSocket(
  //       `ws://localhost:8080/join?roomID=${params.roomID}`
  //     );

  //     webSocketRef.current.addEventListener("open", () => {
  //       webSocketRef.current.send(JSON.stringify({ join: true }));
  //     });

  //     webSocketRef.current.addEventListener("message", async (e) => {
  //       const message = JSON.parse(e.data);

  //       if (message.join) {
  //         callUser();
  //       } else if (message.offer) {
  //         handleOffer(message.offer);
  //       } else if (message.answer) {
  //         console.log("Receiving answer", message.answer);
  //         peerRef.current.setRemoteDescription(
  //           new RTCSessionDescription(message.answer)
  //         );
  //       } else if (message.iceCandidate) {
  //         console.log("Receiving and adding ICE candidates");

  //         try {
  //           await peerRef.current.addIceCandidate(message.iceCandidate);
  //         } catch (err) {
  //           console.log("Error receiving ice candidates", err);
  //         }
  //       }
  //     });
  //   });
  // });

  const initFunc = () => {
    openCamera().then((stream) => {
      userVideo.current.srcObject = stream;
      userStream.current = stream;

      webSocketRef.current = new WebSocket(
        `ws://localhost:8080/join?roomID=${params.roomID}`
      );

      webSocketRef.current.addEventListener("open", () => {
        webSocketRef.current.send(JSON.stringify({ join: true }));
      });

      webSocketRef.current.addEventListener("message", async (e) => {
        const message = JSON.parse(e.data);

        if (message.join) {
          callUser();
        } else if (message.offer) {
          handleOffer(message.offer);
        } else if (message.answer) {
          console.log("Receiving answer", message.answer);
          peerRef.current.setRemoteDescription(
            new RTCSessionDescription(message.answer)
          );
        } else if (message.iceCandidate) {
          console.log("Receiving and adding ICE candidates");

          try {
            await peerRef.current.addIceCandidate(message.iceCandidate);
          } catch (err) {
            console.log("Error receiving ice candidates", err);
          }
        }
      });
    });
  };

  initFunc();

  const handleOffer = async (offer) => {
    console.log("Received offer, creating answer");

    peerRef.current = createPeer();

    await peerRef.current.setRemoteDescription(
      new RTCSessionDescription(offer)
    );

    userStream.current.getTracks().forEach((track) => {
      peerRef.current.addTrack(track, userStream.current);
    });

    const answer = await peerRef.current.createAnswer();
    await peerRef.current.setLocalDescription(answer);

    webSocketRef.current.send(
      JSON.stringify({ answer: peerRef.current.localDescription })
    );
  };

  const callUser = () => {
    console.log("Calling other user");
    peerRef.current = createPeer();

    userStream.current.getTracks().forEach((track) => {
      peerRef.current.addTrack(track, userStream.current);
    });
  };

  const createPeer = () => {
    console.log("Creating peer");

    const peer = new RTCPeerConnection({
      iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
    });

    peer.onnegotiationneeded = handleNegotiationNeeded;
    peer.onicecandidate = handleIceCandidateEvent;
    peer.ontrack = handleTrackEvent;

    return peer;
  };

  const handleNegotiationNeeded = async () => {
    console.log("Creating offer");

    try {
      const myOffer = await peerRef.current.createOffer();
      await peerRef.current.setLocalDescription(myOffer);

      webSocketRef.current.send(
        JSON.stringify({ offer: peerRef.current.localDescription })
      );
    } catch (err) {
      console.log("Error in handleNegotiationNeeded", err);
    }
  };

  const handleIceCandidateEvent = (e) => {
    console.log("Found ICE candidates");

    if (e.candidate) {
      console.log("Candidate", e.candidate);
      webSocketRef.current.send(JSON.stringify({ iceCandidate: e.candidate }));
    }
  };

  const handleTrackEvent = (e) => {
    console.log("Received tracks");
    partnerVideo.current.srcObject = e.streams[0];
  };

  return (
    <div>
      <video autoPlay controls={true} ref={userVideo}></video>
      <video autoPlay controls={true} ref={partnerVideo}></video>
    </div>
  );
};

export default Room;
