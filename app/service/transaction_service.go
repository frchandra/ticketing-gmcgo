package service

import (
	"errors"
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type TransactionService struct {
	txRepo   *repository.TransactionRepository
	userRepo *repository.UserRepository
	seatRepo *repository.SeatRepository
}

func NewTransactionService(txRepo *repository.TransactionRepository, userRepo *repository.UserRepository, seatRepo *repository.SeatRepository) *TransactionService {
	return &TransactionService{txRepo: txRepo, userRepo: userRepo, seatRepo: seatRepo}
}

func (s *TransactionService) CreateTx(userId uint64, seatIds []uint) error {
	txId := uuid.New().String()
	/*	var user model.User
		if result := s.userRepo.GetById(userId, &user); result.Error != nil {
			return result.Error
		}*/

	for _, seatId := range seatIds {
		/*		var seat model.Seat
				if result := s.seatRepo.GetSeatById(&seat, seatId); result.Error != nil {
					return result.Error
				}*/
		newTx := model.Transaction{
			OrderId: txId,
			UserId:  userId,
			SeatId:  seatId,
			//User:         user,
			//Seat:         seat,
			Vendor:       "#",
			Confirmation: "reserved",
		}
		//delete previous failed reservation
		s.txRepo.SoftDeleteTransaction(seatId, userId)

		if result := s.txRepo.InsertOne(&newTx); result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (s *TransactionService) SoftDeleteTransaction() {

}

func (s *TransactionService) SeatsBelongsToUserId(userId uint64) ([]model.Seat, error) {
	var transactions []model.Transaction
	var seats []model.Seat
	if result := s.txRepo.GetLastTxByUserId(&transactions, userId); result.RowsAffected < 1 {
		return seats, errors.New("user belum melakukan pemesanan/transaksi")
	}

	for _, tx := range transactions {
		var seat model.Seat
		s.seatRepo.GetSeatById(&seat, tx.SeatId)
		if tx.Confirmation == "reserved" {
			seat.Status = "reserved_by_me"
		}
		if tx.Confirmation == "settlement" {
			seat.Status = "purchased_by_me"
		}
		seats = append(seats, seat)
	}
	return seats, nil
}

func (s *TransactionService) GetUserTransactionDetails(userId uint64) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if result := s.txRepo.GetUserTransactionDetails(&transactions, userId); result.Error != nil {
		return transactions, result.Error
	}
	return transactions, nil

}

func (s *TransactionService) PrepareTransactionData(userId uint64) snap.Request {
	txDetails, _ := s.GetUserTransactionDetails(userId)
	var grossAmt int64
	var itemDetails []midtrans.ItemDetails

	customerDetails := midtrans.CustomerDetails{
		FName: txDetails[0].User.Name,
		LName: "",
		Email: txDetails[0].User.Email,
		Phone: txDetails[0].User.Phone,
	}

	for _, tx := range txDetails {
		grossAmt += int64(tx.Seat.Price)
		itemDetail := midtrans.ItemDetails{
			ID:    string(tx.SeatId),
			Price: int64(tx.Seat.Price),
			Qty:   1,
			Name:  tx.Seat.Name,
		}
		itemDetails = append(itemDetails, itemDetail)
	}

	var snapRequest snap.Request = snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  txDetails[0].OrderId,
			GrossAmt: grossAmt,
		},
		CustomerDetail: &customerDetails,
		Items:          &itemDetails,
	}
	return snapRequest
}
