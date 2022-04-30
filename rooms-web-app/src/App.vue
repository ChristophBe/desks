<template>
  <v-app>
    <v-app-bar>
      <v-card-title>Rooms</v-card-title>
      <v-spacer></v-spacer>
      <v-avatar
          v-if="user"
          color="primary"
          size="48"
      >
        <span class="white--text text-h5">{{user.firstname.charAt(0).toUpperCase() + user.lastname.charAt(0).toUpperCase()}}</span>
      </v-avatar>
    </v-app-bar>
    <v-main>


      <v-container>
        <FakeLogin v-if="!user && !loading"/>

        <v-card v-else-if="loading">
          <v-card-title>Wait for it</v-card-title>
        </v-card>
        <v-card v-else>
          <v-card-title>You are logged in!</v-card-title>
        </v-card>
      </v-container>


    </v-main>
  </v-app>
</template>

<script>

import FakeLogin from "@/components/GetUser";
import { mapState } from 'vuex'

export default {
  username: 'App',

  components: {
    FakeLogin,
  },

  mounted() {
    this.$store.dispatch("fetchCurrentUser")
  },
  computed: mapState(['user',"loading"])
}
</script>
