<template>
  <v-dialog
      v-model="showBookingFormDialog"
      persistent
  >
    <BookingForm @close="showBookingFormDialog = false" :booking="bookingToEdit"></BookingForm>
  </v-dialog>
  <v-card v-if="todaysBookings.length > 0" class="mb-3">
    <v-card-title>Today</v-card-title>

    <v-expansion-panels>
      <v-expansion-panel v-for="booking in todaysBookings" :key="booking.id">
        <v-expansion-panel-title class="pl-4">
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
          @click="openCreateBookingDialog()"
      >
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </v-card-title>

    <v-alert v-if="upcomingBookings.length  <= 0">
      You have currently no upcoming desk bookings.
    </v-alert>
    <BookingsTable v-else :bookings="upcomingBookings" @editBooking="(booking) => openEditeBookingDialog(booking)"></BookingsTable>
  </v-card>
</template>

<script lang="ts">

import {defineComponent} from "vue";
import {mapActions, mapGetters} from "vuex";
import moment from "moment";
import BookingsTable from "@/components/booking-components/BookingsTable.vue";
import AlsoInTheRoom from "@/components/booking-components/AlsoInTheRoom.vue";
import Booking from "@/models/Booking";
import BookingForm from "@/components/booking-components/BookingForm.vue";

interface bookingViewData{
  showBookingFormDialog:boolean
  bookingToEdit:Booking|null
}

export default defineComponent({
  name: "BookingsView",
  components: {AlsoInTheRoom, BookingsTable, BookingForm},
  data: ():bookingViewData  => ({
    showBookingFormDialog: false,
    bookingToEdit: null
  }),
  computed: mapGetters('bookings', ['upcomingBookings', 'todaysBookings']),
  mounted() {
    this.fetchBookings()
  },
  methods: {
    ...mapActions("bookings",["fetchBookings"]),
    formatDate(date:Date) {
      return moment(date).format("DD.MM.YYYY")
    },
    formatTime(date:Date) {
      return moment(date).format("HH:mm")
    },
    openCreateBookingDialog(){
      this.showBookingFormDialog = true
      this.bookingToEdit = null;
    },
    openEditeBookingDialog(booking:Booking){
      this.showBookingFormDialog = true
      this.bookingToEdit = booking
    }
  }

})

</script>

<style scoped>

</style>
