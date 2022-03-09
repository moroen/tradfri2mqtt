<template>
  <q-list>
    <q-item link v-for="page in navItems" :key="page.name" :to="page.destination">
      <q-item-section avatar>
        <q-icon :name="page.avatar" />
      </q-item-section>

      <q-item-section>
        <q-item-label>{{ page.name }}</q-item-label>
        <q-item-label caption>{{ page.caption }}</q-item-label>
      </q-item-section>
    </q-item>
  </q-list>
</template>

<script setup>
const navItems = [
  {
    name: "Tradfri",
    destination: "/tradfri",
    caption: "Manage IKEA-Gateway and devices",
    avatar: "lightbulb",
  },
  {
    name: "MQTT",
    destination: "/mqtt",
    caption: "MQTT settings and tools",
    avatar: "share",
  },
  {
    name: "Interface",
    destination: "/interface",
    caption: "Configure web interface",
    avatar: "settings",
  },
  {
    name: "Log",
    destination: "/log",
    caption: "View log",
    avatar: "wysiwyg",
  },
  {
    name: "About",
    destination: "/about",
    caption: "",
    avatar: "info",
  },
];
</script>
