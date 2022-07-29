<template>
  <v-app>
    <v-app-bar>
      <v-card-title>Desks</v-card-title>
      <v-spacer></v-spacer>
      <v-avatar
          v-if="user"
          color="primary"
          size="48"
      >
        <span
            class="white--text text-h5">{{
            user.givenName.charAt(0).toUpperCase() + user.familyName.charAt(0).toUpperCase()
          }}</span>
      </v-avatar>
    </v-app-bar>
    <v-main>


      <v-container>
        <Login v-if="!user && !loading"/>

        <v-card v-else-if="loading">
          <v-card-title>Wait for it</v-card-title>
        </v-card>
        <BookingsView v-else></BookingsView>
      </v-container>

      <app-notifications></app-notifications>
    </v-main>

  </v-app>
</template>

<script lang="ts">
import {Options, Vue} from 'vue-class-component';
import Login from './components/views/LoginView.vue';
import {mapState} from "vuex";
import Bookings from "./components/views/BookingsView.vue";
import BookingsView from "@/components/views/BookingsView.vue";
import Notifications from "@/components/Notifications.vue";
import AppNotifications from "@/components/Notifications.vue";

@Options({
  components: {
    AppNotifications,
    Notifications,
    BookingsView,
    Bookings,
    Login: Login,
  },

  mounted() {
    this.$store.dispatch("user/fetchCurrentUser")
  },
  computed: mapState('user', ['user', 'loading'])
})
export default class App extends Vue {
}
</script>


