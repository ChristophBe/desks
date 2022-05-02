<template>
  <v-dialog
      v-model="dialog"
  >
    <AddBookingForm @close="dialog=false"></AddBookingForm>
  </v-dialog>
  
  <v-card>
    <v-card-title>
      <div>
        My Desk Bookings
      </div>
      <v-spacer></v-spacer>
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


export default defineComponent({
  name: "BookingsView",
  components: {BookingsTable, AddBookingForm},
  data: () => ({
    dialog: false
  }),
  computed: mapGetters('bookings', ['upcomingBookings']),
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