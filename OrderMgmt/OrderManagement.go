package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

type PO struct {
	Order_Id       string `json:"order_Id"`
	Order_Desc     string `json:"order_desc"`
	Order_Quantity string `json:"order_quantity"`
	Assigned_To_Id string `json:"assigned_to_id"`
	Created_By_Id  string `json:"created_by_id"`
	Order_Status   string `json:"order_status"`
	Asset_ID       string `json:"asset_ID"`
}

type SimpleChaincode struct {
}

func main() {

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)

	}
}

func (t *SimpleChaincode) convert(row shim.Row) PO {
	var po = PO{}

	po.Order_Id = row.Columns[0].GetString_()
	po.Order_Desc = row.Columns[1].GetString_()
	po.Order_Quantity = row.Columns[2].GetString_()
	po.Assigned_To_Id = row.Columns[3].GetString_()
	po.Created_By_Id = row.Columns[4].GetString_()
	po.Order_Status = row.Columns[5].GetString_()
	po.Asset_ID = row.Columns[6].GetString_()

	return po
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	err = stub.CreateTable("OEM", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{"OrderId", shim.ColumnDefinition_STRING, true},
		&shim.ColumnDefinition{"Order_Desc", shim.ColumnDefinition_STRING, false},
		&shim.ColumnDefinition{"Order_Quantity", shim.ColumnDefinition_STRING, false},
		&shim.ColumnDefinition{"Assigned_To_Id", shim.ColumnDefinition_STRING, false},
		&shim.ColumnDefinition{"Created_By_Id", shim.ColumnDefinition_STRING, false},
		&shim.ColumnDefinition{"Order_Status", shim.ColumnDefinition_STRING, false},
		&shim.ColumnDefinition{"Asset_ID", shim.ColumnDefinition_STRING, false}})

	if err != nil {
		return nil, errors.New("OEM table not created")
	}

	err = stub.CreateTable("TIER1", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{"subOrderId", shim.ColumnDefinition_STRING, true},
		&shim.ColumnDefinition{"OrderId", shim.ColumnDefinition_STRING, true},
		&shim.ColumnDefinition{"Order_Desc", shim.ColumnDefinition_STRING, false},
		&shim.ColumnDefinition{"Order_Quantity", shim.ColumnDefinition_STRING, false},
		&shim.ColumnDefinition{"Assigned_To_Id", shim.ColumnDefinition_STRING, false},
		&shim.ColumnDefinition{"Created_By_Id", shim.ColumnDefinition_STRING, false},
		&shim.ColumnDefinition{"Order_Status", shim.ColumnDefinition_STRING, false},
		&shim.ColumnDefinition{"Asset_ID", shim.ColumnDefinition_STRING, false}})

	if err != nil {
		return nil, errors.New("TIER1 table not created")
	}

	orderId := "1000"
	subOrderId := "12345"

	stub.PutState("orderId", []byte(orderId))
	stub.PutState("subOrderId", []byte(subOrderId))

	return nil, nil

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "createOrder" {

		return createOrder(stub, args)
	}
	return nil, nil
}
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "fetchAllOrders" {

		return fetchAllOrders(stub, args)
	}
	return nil, nil

}

func createOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//OrderId
	byteOrderId, err := stub.GetState("orderId")
	strOrderId := string(byteOrderId)
	intOrderId, _ := strconv.Atoi(strOrderId)

	currentId := intOrderId + 1
	str := strconv.Itoa(currentId)
	strCurrentId := "PO" + strconv.Itoa(currentId)
	stub.PutState("orderId", []byte(str))

	fmt.Println(strCurrentId)
	fmt.Println(args[0])
	fmt.Println(args[1])
	fmt.Println(args[3])

	col1Val := strCurrentId
	col2Val := args[0]
	col3Val := args[1]
	col4Val := "Tier 1"
	col5Val := "OEM"
	col6Val := "Order created. Pending with Tier1"
	col7Val := args[3]

	var columns []*shim.Column

	col0 := shim.Column{Value: &shim.Column_String_{String_: col1Val}}
	col1 := shim.Column{Value: &shim.Column_String_{String_: col2Val}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: col3Val}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: col4Val}}
	col4 := shim.Column{Value: &shim.Column_String_{String_: col5Val}}
	col5 := shim.Column{Value: &shim.Column_String_{String_: col6Val}}
	col6 := shim.Column{Value: &shim.Column_String_{String_: col7Val}}

	columns = append(columns, &col0)
	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)
	columns = append(columns, &col5)
	columns = append(columns, &col6)

	row := shim.Row{Columns: columns}
	ok, err := stub.InsertRow("OEM", row)

	if err != nil {
		return nil, fmt.Errorf("insertTableOne operation failed. %s", err)
		panic(err)

	}
	if !ok {
		return []byte("Row with given key" + args[0] + " already exists"), errors.New("insertTableOne operation failed. Row with given key already exists")
	}
	return nil, nil
}

func fetchAllOrders(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var columns []shim.Column
	rowChannel, err := stub.GetRows("OEM", columns)

	orderArray := []PO{}

	for {
		select {

		case row, ok := <-rowChannel:

			if !ok {
				rowChannel = nil
			} else {

				fmt.Println("Inside Else of for loop in query")
				po := PO{}

				po.Order_Id = row.Columns[0].GetString_()
				po.Order_Desc = row.Columns[1].GetString_()
				po.Order_Quantity = row.Columns[2].GetString_()
				po.Assigned_To_Id = row.Columns[3].GetString_()
				po.Created_By_Id = row.Columns[4].GetString_()
				po.Order_Status = row.Columns[5].GetString_()
				po.Asset_ID = row.Columns[6].GetString_()

				orderArray = append(orderArray, po)
			}

		}
		if rowChannel == nil {
			break
		}
	}

	jsonRows, err := json.Marshal(orderArray)

	if err != nil {
		return nil, fmt.Errorf("getRowsTableFour operation failed. Error marshaling JSON: %s", err)
	}

	return jsonRows, nil

}
