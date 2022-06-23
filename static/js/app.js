const Prompt = () => {
  const toast = (c) => {
    const { title = "", icon = "success", position = "top-end" } = c;

    const Toast = Swal.mixin({
      toast: true,
      position: position,
      title: title,
      icon: icon,
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.addEventListener("mouseenter", Swal.stopTimer);
        toast.addEventListener("mouseleave", Swal.resumeTimer);
      },
    });

    Toast.fire({});
  };

  const success = (c) => {
    const { title = "", text = "", footer = "" } = c;
    Swal.fire({
      icon: "success",
      title: title,
      text: text,
      footer: footer,
    });
  };

  const error = (c) => {
    const { title = "", text = "", footer = "" } = c;
    Swal.fire({
      icon: "error",
      title: title,
      text: text,
      footer: footer,
    });
  };

  const custom = async (c) => {
    const { msg = "", title = "", icon = "", showConfirmButton = true } = c;

    const { value: result } = await Swal.fire({
      title: title,
      html: msg,
      icon: icon,
      showConfirmButton: showConfirmButton,
      backdrop: false,
      focusConfirm: false,
      showCancelButton: true,
      willOpen: () => {
        if (c.willOpen !== undefined) {
          c.willOpen();
        }
      },
      preConfirm: () => {
        return [
          document.getElementById("start").value,
          document.getElementById("end").value,
        ];
      },
      didOpen: () => {
        if (c.didOpen !== undefined) {
          c.didOpen();
        }
      },
    });

    if (result) {
      if (result.dismiss !== Swal.DismissReason.cancel) {
        if (result.value !== "") {
          if (c.callback !== undefined) {
            c.callback(result);
          }
        } else {
          c.callback(false);
        }
      } else {
        c.callback(false);
      }
    }
  };

  return {
    toast: toast,
    success: success,
    error: error,
    custom: custom,
  };
};

// const displayDates = (roomID) => {
//   document
//     .getElementById("check-availability-button")
//     .addEventListener("click", () => {
//       let html = `
//           <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
//             <div class="row">
//               <div class="col">
//                 <div class="row" id="reservation-dates-modal">
//                   <div class="col">
//                     <input disabled required class="form-control" autocomplete="off" type="text" name="start" id="start" placeholder="Arrival"/>
//                   </div>
//                   <div class="col">
//                     <input disabled required class="form-control" autocomplete="off" type="text" name="end" id="end" placeholder="Departure"/>
//                   </div>
//                 </div>
//               </div>
//             </div>
//           </form>
//         `;

//       attention.custom({
//         msg: html,
//         title: "Choose your date",
//         willOpen: () => {
//           const elem = document.getElementById("reservation-dates-modal");
//           const rp = new DateRangePicker(elem, {
//             format: "yyyy-mm-dd",
//             showOnFocus: true,
//             minDate: new Date(),
//           });
//         },
//         didOpen: () => {
//           document.getElementById("start").removeAttribute("disabled");
//           document.getElementById("end").removeAttribute("disabled");
//         },
//         callback: (result) => {
//           console.log("called");
//           let form = document.getElementById("check-availability-form");
//           let formData = new FormData(form);
//           formData.append("csrf_token", "{{.CSRFToken}}");
//           formData.append("room_id", roomID);
//           fetch("/search-availability-json", { method: "post", body: formData })
//             .then((res) => res.json())
//             .then((data) => {
//               if (data.ok) {
//                 attention.custom({
//                   icon: "success",
//                   msg:
//                     '<p>Room is available</p><p><a href="/book-room?id=' +
//                     data.roomID +
//                     "&s=" +
//                     data.startDate +
//                     "&e=" +
//                     data.endDate +
//                     '" class="btn btn-primary">Book Now</a></p>',
//                   showConfirmButton: false,
//                 });
//               } else {
//                 attention.error({ msg: "No Availability" });
//               }
//             });
//         },
//       });
//     });
// };
