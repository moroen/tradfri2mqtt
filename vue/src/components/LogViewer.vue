<template>
  <div>
    <q-checkbox v-model="showlog" @update:model-value="onLogToggle">Show log</q-checkbox>
  </div>

  <q-scroll-area
    class="col q-pa-sm"
    :thumb-style="thumbStyle"
    :content-style="contentStyle"
    :content-active-style="contentActiveStyle"
    style="col"
  >
    <div v-for="(message, idx) in log" :key="log">
      <span>{{ message.time }}</span>
      <span>{{ message.level }}</span>
      <span>{{ message.message }}</span>
    </div>
  </q-scroll-area>
</template>

<script setup>
import { ref } from "vue";
import { mapFields } from "vuex-map-fields";
import { useStore } from "vuex";

// const connection = new WebSocket(`ws://${window.location.host}/api/ws`);
const messages = ref([]);
var connection;
const showlog = ref(false);

const logtext = ref("");

const store = useStore();

const onLogToggle = function (value) {
  console.log("Hi");
  var cmd;

  if (value) {
    cmd = "start";
  } else {
    cmd = "stop";
  }

  store.dispatch(
    "sendWSCommand",
    JSON.stringify({
      class: "log",
      command: cmd,
    })
  );

  if (!value) {
    store.commit("clearLog");
  }
  // if (value) {
  //   connection = new WebSocket(`ws://${window.location.host}/api/ws`);
  //   connection.onopen = function (event) {
  //     console.log(event);
  //     console.log("Successfully connected to the echo websocket server...");
  //   };

  //   connection.onmessage = (event) => {
  //     // console.log("Message received");
  //     var obj = JSON.parse(event.data);
  //     console.log(obj);
  //     messages.value.push(obj);
  //     console.log(messages.value[0]);
  //     logtext.value = event.data;
  //     // console.log(messages);
  //   };
  // } else {
  //   messages.value = [];
  //   connection.close();
  // }
};

const contentStyle = {
  backgroundColor: "rgba(0,0,0,0.02)",
  color: "#555",
};

const contentActiveStyle = {
  backgroundColor: "#eee",
  color: "black",
};
const thumbStyle = {
  right: "2px",
  borderRadius: "5px",
  backgroundColor: "#027be3",
  width: "5px",
  opacity: 0.75,
};

const SendMessage = function () {
  console.log("OnClick");
  connection.send("test data");
};

const styleFn = function (offset, height) {
  let pageheight = height - offset;
  console.log("PageHeight: " + pageheight);
  return "height: " + pageheight + "px";
};
</script>

<script>
export default {
  computed: {
    ...mapFields(["websocket.log"]),
  },
};
</script>
