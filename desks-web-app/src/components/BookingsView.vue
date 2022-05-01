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
    <v-table>
      <thead>
      <tr>
        <th class="text-left">
          Room
        </th>
        <th class="text-left">
          Date
        </th>
        <th class="text-left">
          Start
        </th>
        <th class="text-left">
          End
        </th>
      </tr>
      </thead>
      <tbody>
      <tr v-if="bookings.length  <= 0">
        <td colspan="3">
          <v-alert class="mt-2 mb-2">
            You have currently no upcoming desk bookings.
          </v-alert>
        </td>
      </tr>
      <tr v-else
          v-for="item in bookings"
          :key="item.name"
      >
        <td>{{ item.room.name }}</td>
        <td>{{ formatDate(item.start)}}</td>
        <td>{{ formatTime(item.start) }}</td>
        <td>{{ formatTime(item.end) }}</td>
      </tr>
      </tbody>
    </v-table>
  </v-card>
</template>

<script>
import AddBookingForm from "@/components/AddBookingForm";
import {defineComponent} from "vue";
import {mapState} from "vuex";
import moment from "moment";


export default defineComponent({
  name:"BookingsView",
  components: {AddBookingForm},
  data:()=>({
    dialog: false
  }),
  computed: mapState(['bookings']),
  mounted() {
    this.$store.dispatch("fetchBookings")
  },
  methods:{
    formatDate(date){
      return moment(date).format("DD.MM.YYYY")
    },
    formatTime(date){
      return moment(date).format("HH:mm")
    }
  }

})

</script>

<style scoped>

</style>