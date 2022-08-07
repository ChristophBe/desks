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
        {{
          user.givenName.charAt(0).toUpperCase() + user.familyName.charAt(0).toUpperCase()
        }}
      </v-avatar>
    </v-app-bar>
    <v-main>

      <v-fade-transition>
        <app-loading v-if="!user || loading" text="Authenticating"></app-loading>


        <BookingsView v-else></BookingsView>
      </v-fade-transition>


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
import AppLoading from "@/components/Loading.vue";

@Options({
  components: {
    AppLoading,
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


