<template>
  <v-dialog
      v-model="showBookingFormDialog"
      persistent
      width="400"
  >
    <BookingFormDialogue @close="showBookingFormDialog = false" :booking="bookingToEdit"></BookingFormDialogue>
  </v-dialog>
  <v-dialog
      v-model="showBookingDialog"
      width="400"
  >
    <booking-details @close="showBookingDialog = false" :booking="bookingToShow"
                     @edit="openEditeBookingDialog(bookingToShow)"></booking-details>
  </v-dialog>
  <v-fade-transition>

    <loading v-if="(loading && !bookingsFetched) || (bookingDefaultsLoading && !bookingDefaultsFetched)"></loading>
    <v-container v-else>
      <desk-availabilty @book="openEditeBookingDialog"/>

      <v-row class="mt-8">
        <v-col>

          <h1 class="text-h4">My Desk Bookings</h1>
        </v-col>

        <v-col align-self="center" class="d-flex justify-end">
          <v-btn
              variant="flat"
              class="rounded-pill"
              @click="openCreateBookingDialog()"
          >
            <template v-slot:prepend>

              <v-icon>mdi-plus</v-icon>
            </template>
            Add Booking
          </v-btn>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
            <v-card-text v-if="myUpcomingBookings.length <= 0">
              <v-alert>
                You have no upcoming desk bookings.
              </v-alert>
            </v-card-text>
            <BookingsCards
                v-else
                :bookings="myUpcomingBookings"
                @editBooking="(booking) => openEditeBookingDialog(booking)"
                @openBooking="(booking) => openShowBookingDialog(booking)"
            ></BookingsCards>
        </v-col>
      </v-row>


    </v-container>
  </v-fade-transition>


</template>

<script lang="ts">

import {defineComponent} from "vue";
import {mapActions, mapGetters, mapState} from "vuex";
import BookingsCards from "@/components/booking-components/BookingsCards.vue";
import Booking from "@/models/Booking";
import BookingDetails from "@/components/booking-components/BookingDetails.vue";
import BookingFormDialogue from "@/components/booking-components/BookingFormDialogue.vue";
import DeskAvailabilty from "@/components/booking-components/WeekOverviewRow.vue";
import Loading from "@/components/Loading.vue";
import moment from "moment";
import BookingUtils from "@/utils/booking-utils";


interface bookingViewData {
  showBookingFormDialog: boolean
  bookingToEdit: Booking | null
  showBookingDialog: boolean
  bookingToShow: Booking | null
  show: boolean
}

export default defineComponent({
  name: "BookingsView",
  components: {DeskAvailabilty, BookingFormDialogue, BookingDetails, BookingsCards, Loading},
  data: (): bookingViewData => ({
    showBookingFormDialog: false,
    bookingToEdit: null,
    showBookingDialog: false,
    bookingToShow: null,
    show: false
  }),
  computed: {
    ...mapState('bookings', ['loading', "bookingsFetched"]),
    ...mapState('defaults', ['bookingDefaultsLoading', 'bookingDefaultsFetched']),
    ...mapGetters('bookings', ['myUpcomingBookings', 'myBookingsOfTheDay']),
  },
  mounted() {
    this.fetchMyBookings()
    this.fetchBookingDefaults()
  },
  methods: {
    ...mapActions("bookings", ["fetchMyBookings", "leaveEarly"]),
    ...mapActions("defaults", ["fetchBookingDefaults"]),
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
    },
    onLeaveEarly(booking: Partial<Booking>) {

      this.leaveEarly(booking)
    },
    hasBookingsForToday() {
      return this.myBookingsOfTheDay.length > 0
    },
    isOngoing(booking: Booking) {
      return moment(booking.start).isSameOrBefore(moment.now()) && moment(booking.end).isSameOrAfter(moment.now())
          && !BookingUtils.roundToNextMinute().isSame(booking.end)
    },
  }

})

</script>

