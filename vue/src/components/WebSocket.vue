<template>
  <div>
    <p v-for="(message, idx) in messages" :key="idx">{{ message }}</p>
  </div>
</template>

<script>
// We store the reference to the SSE client out here

import { toHandlers } from "vue";

// so we can access it from other methods
let sseClient;

let count = 0;

export default {

  data() {
    return {
      messages: [],
      connection: null,
    };
  },
  mounted() {
    this.connection = new WebSocket("ws://localhost:8321/api/ws")
    this.connection.onmessage = (event) => {
      console.log(event)
    }

    this.connection.onopen = function (event) {
      console.log(event)
      console.log("Successfully connected to the echo websocket server...")
    }
  }
}
</script>