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
  <v-fade-transition>

    <loading v-if="(loading && !bookingsFetched) || (bookingDefaultsLoading && !bookingDefaultsFetched)"></loading>
    <v-container v-else>

      <v-row class="mb-3">
        <v-col :md="hasBookingsForToday() ? 6 : 12" cols="12">
          <desk-availabilty @book="openEditeBookingDialog"/>
        </v-col>
        <v-col cols="12" md="6" v-if="hasBookingsForToday()" v-bind:id="today">
          <v-card>
            <v-card-title>Today</v-card-title>

            <v-list>
              <v-list-item
                  v-for="booking in myBookingsOfTheDay"
                  :key="booking.id"
                  @click="openShowBookingDialog(booking)"
              >
                <v-list-item-header>
                  <v-list-item-title>{{ booking.room.name }}</v-list-item-title>
                  <v-list-item-subtitle>{{ $format.timeRange(booking.start, booking.end) }}</v-list-item-subtitle>

                </v-list-item-header>
                <v-list-item-action v-if="isOngoing(booking)"><v-btn @click.stop.prevent="onLeaveEarly(booking) " variant="text" >leave now</v-btn></v-list-item-action>
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
            <v-card-text v-if="myUpcomingBookings.length <= 0">
              <v-alert>
                You have no upcoming desk bookings.
              </v-alert>
            </v-card-text>
            <BookingsTable
                v-else
                :bookings="myUpcomingBookings"
                @editBooking="(booking) => openEditeBookingDialog(booking)"
                @openBooking="(booking) => openShowBookingDialog(booking)"
            ></BookingsTable>
          </v-card>
        </v-col>
      </v-row>


    </v-container>
  </v-fade-transition>


</template>

<style scoped>
  .v-row .v-card {
    height: 100%;
  }
  
  .v-card button.v-btn {
    border: 1px solid rgba(0,0,0, 0.06);
  }
  
  .v-card button.v-btn.v-theme--dark {
    border: 1px solid rgba(255,255,255, 0.06);
  }
  
  .v-card button.v-btn {
    padding: 0.5em 1.25em;
  }
</style>

<script lang="ts">

import {defineComponent} from "vue";
import {mapActions, mapGetters, mapState} from "vuex";
import BookingsTable from "@/components/booking-components/BookingsTable.vue";
import Booking from "@/models/Booking";
import BookingDetails from "@/components/booking-components/BookingDetails.vue";
import BookingFormDialogue from "@/components/booking-components/BookingFormDialogue.vue";
import DeskAvailabilty from "@/components/booking-components/DeskAvailabilityCard.vue";
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
  components: {DeskAvailabilty, BookingFormDialogue, BookingDetails, BookingsTable, Loading},
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
    ...mapActions("bookings", ["fetchMyBookings","leaveEarly"]),
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
    isOngoing(booking: Booking){
      return moment(booking.start).isSameOrBefore(moment.now()) && moment(booking.end).isSameOrAfter(moment.now())
      && !BookingUtils.roundToNextMinute().isSame(booking.end)
    },
  }

})

</script>

