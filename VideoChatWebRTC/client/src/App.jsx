import React from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";

import CreateRoom from "./components/CreateRoom";
import Room from "./components/Room";

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Routes>
          <Route path="/" exact Component={CreateRoom}></Route>
          <Route path="/room/:roomID" Component={Room}></Route>
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
