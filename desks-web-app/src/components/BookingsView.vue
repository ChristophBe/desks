<template>
  <v-dialog
      v-model="dialog"
      persistent
  >
    <AddBookingForm @close="dialog=false"></AddBookingForm>
  </v-dialog>
  <v-card v-if="todaysBookings.length > 0" class="mb-3">
    <v-card-title>Today</v-card-title>

    <v-expansion-panels>
      <v-expansion-panel v-for="booking in todaysBookings" :key="booking.id">
        <v-expansion-panel-title class="pl-4 pr-8">
          <span>{{ booking.room.name }}</span>
          <v-spacer></v-spacer>
          <span class="mr-4">
            {{ formatTime(booking.start) }} - {{ formatTime(booking.end) }}
          </span>

        </v-expansion-panel-title>
        <v-expansion-panel-text class="px-0">
          <AlsoInTheRoom :room-id="booking.room.id" :date="booking.start"></AlsoInTheRoom>
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>
  </v-card>
  <v-card>
    <v-card-title class="d-flex">

      <div class="mr-auto">
        My Desk Bookings
      </div>

      <v-btn
          icon
          elevation="0"
          @click="dialog=true"
      >
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </v-card-title>

    <v-alert v-if="upcomingBookings.length  <= 0">
      You have currently no upcoming desk bookings.
    </v-alert>
    <BookingsTable v-else :bookings="upcomingBookings"></BookingsTable>
  </v-card>
</template>

<script>
import AddBookingForm from "@/components/AddBookingForm";
import {defineComponent} from "vue";
import {mapGetters} from "vuex";
import moment from "moment";
import BookingsTable from "@/components/BookingsTable";
import AlsoInTheRoom from "@/components/AlsoInTheRoom";


export default defineComponent({
  name: "BookingsView",
  components: {AlsoInTheRoom, BookingsTable, AddBookingForm},
  data: () => ({
    dialog: false
  }),
  computed: mapGetters('bookings', ['upcomingBookings', 'todaysBookings']),
  mounted() {
    this.$store.dispatch("bookings/fetchBookings")
  },
  methods: {
    formatDate(date) {
      return moment(date).format("DD.MM.YYYY")
    },
    formatTime(date) {
      return moment(date).format("HH:mm")
    }
  }

})

</script>

<style scoped>

</style>