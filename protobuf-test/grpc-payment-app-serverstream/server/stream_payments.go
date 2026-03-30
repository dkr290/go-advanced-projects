package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/dkr290/go-advanced-projects/protobuf-test/grpc-payment-app-serverstream/proto"
	"google.golang.org/grpc"
)


func(s *Server) Payment(req *proto.PaymentRequest, stream grpc.ServerStreamingServer[proto.PaymentResponse]) error  {

  userId := req.UserId
  log.Println("User subscribed to the payment stream ", userId)
	ctx := stream.Context()
  for i := range 5 {
		log.Println("Number of payment",i)
		select {
		 case <- ctx.Done():
			 log.Printf("User %s is disconnected", userId)
			 return nil
	  default:
			payment := &proto.PaymentResponse{
         PaymentId: fmt.Sprintf("pay_%d",rand.Intn(999999)),
         Amount: randomAmount(),
				 Currency: "US",
				 Status: randomStatus(),
				 Timestamp: time.Now().Unix(),
			}
			if err := stream.Send(payment);err != nil {
				log.Printf("Error sending to %s: %v",userId,err )
			}
			time.Sleep(20* time.Millisecond)
		}
	}
	log.Printf("Finished streaming payments for %s", userId)
	return nil
}




func randomAmount() float64 {
	min := 10.0
	max := 300.0
	r := min + rand.Float64() * (max - min)
	return r
}
func randomStatus() string{
	s := []string{"PENDING","COMPLETED","FAILED"}
	return s[rand.Intn(len(s))]
} 
