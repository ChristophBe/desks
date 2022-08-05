<template>
  <v-dialog
      v-model="showBookingFormDialog"
      persistent
  >
    <BookingFormDialogue @close="showBookingFormDialog = false" :booking="bookingToEdit"></BookingFormDialogue>
  </v-dialog>
  <v-dialog
      v-model="showBookingDialog"
  >
    <booking-details @close="showBookingDialog = false" :booking="bookingToShow"
                     @edit="openEditeBookingDialog(bookingToShow)"></booking-details>
  </v-dialog>
  <v-container>
    <v-row class="mb-3">
      <v-col cols="6" xm="12">
        <v-card v-if="todaysBookings.length > 0">
          <v-card-item>
            <div class="d-flex justify-space-between align-center">
              <div>
                <v-card-title>Desk Availability</v-card-title>
                <v-card-subtitle>HotSprings Office</v-card-subtitle>
              </div>
              <v-chip color="success">1/7 desks booked today</v-chip>
            </div>
          </v-card-item>

          <v-card-actions>
            <v-btn variant="text">
              Book desk
            </v-btn>
          </v-card-actions>

        </v-card>
      </v-col>
      <v-col cols="6" xm="12">
        <v-card v-if="todaysBookings.length > 0">
          <v-card-title>Today</v-card-title>

          <v-list>
            <v-list-item
                v-for="booking in todaysBookings"
                :key="booking.id"
                @click="openShowBookingDialog(booking)"
            >
              <v-list-item-header>
                <v-list-item-title>{{ booking.room.name }}</v-list-item-title>
                <v-list-item-subtitle>{{ $format.timeRange(booking.start, booking.end) }}</v-list-item-subtitle>
              </v-list-item-header>
            </v-list-item>
          </v-list>
        </v-card>

      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
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

          <v-card-text v-if="upcomingBookings.length <= 0">
            <v-alert>
              You have no upcoming desk bookings.
            </v-alert>
          </v-card-text>
          <BookingsTable
              v-else
              :bookings="upcomingBookings"
              @editBooking="(booking) => openEditeBookingDialog(booking)"
              @openBooking="(booking) => openShowBookingDialog(booking)"
          ></BookingsTable>
        </v-card>
      </v-col>
    </v-row>

  </v-container>


</template>

<script lang="ts">

import {defineComponent} from "vue";
import {mapActions, mapGetters} from "vuex";
import BookingsTable from "@/components/booking-components/BookingsTable.vue";
import AlsoInTheRoom from "@/components/booking-components/AlsoInTheRoom.vue";
import Booking from "@/models/Booking";
import BookingDetails from "@/components/booking-components/BookingDetails.vue";
import BookingFormDialogue from "@/components/booking-components/BookingFormDialogue.vue";


interface bookingViewData {
  showBookingFormDialog: boolean
  bookingToEdit: Booking | null
  showBookingDialog: boolean
  bookingToShow: Booking | null
  show: boolean
}

export default defineComponent({
  name: "BookingsView",
  components: {BookingFormDialogue, BookingDetails, AlsoInTheRoom, BookingsTable},
  data: (): bookingViewData => ({
    showBookingFormDialog: false,
    bookingToEdit: null,
    showBookingDialog: false,
    bookingToShow: null,
    show: false
  }),
  computed: mapGetters('bookings', ['upcomingBookings', 'todaysBookings']),
  mounted() {
    this.fetchBookings()
  },
  methods: {
    ...mapActions("bookings", ["fetchBookings"]),
    openCreateBookingDialog() {
      this.showBookingFormDialog = true
      this.showBookingDialog = false
      this.bookingToEdit = null;
    },
    openEditeBookingDialog(booking: Booking) {
      this.showBookingFormDialog = true
      this.showBookingDialog = false
      this.bookingToEdit = booking
    },
    openShowBookingDialog(booking: Booking) {
      this.showBookingFormDialog = false
      this.showBookingDialog = true
      this.bookingToShow = booking
    }
  }

})

</script>

<style scoped>

</style>
