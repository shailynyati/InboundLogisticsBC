package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

type PO_tier1 struct {
	Order_Id       string `json:"order_Id"`
	Order_Desc     string `json:"order_desc"`
	Order_Quantity string `json:"order_quantity"`
	Assigned_To_Id string `json:"assigned_to_id"`
	Created_By_Id  string `json:"created_by_id"`
	SubOrder_Id    string `json:"subOrder_Id"`
	Order_Status   string `json:"order_status"`
	Asset_ID       string `json:"asset_ID"`
}

type PO_OEM struct {
	Order_Id       string `json:"order_Id"`
	Order_Desc     string `json:"order_desc"`
	Order_Quantity string `json:"order_quantity"`
	Assigned_To_Id string `json:"assigned_to_id"`
	Created_By_Id  string `json:"created_by_id"`
	Order_Status   string `json:"order_status"`
	Asset_ID       string `json:"asset_ID"`
}

func main() {

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)

	}
}

func (t *SimpleChaincode) convert(row shim.Row) PO_tier1 {
	var po = PO_tier1{}

	po.Order_Id = row.Columns[0].GetString_()
	po.Order_Desc = row.Columns[1].GetString_()
	po.Order_Quantity = row.Columns[2].GetString_()
	po.Assigned_To_Id = row.Columns[3].GetString_()
	po.Created_By_Id = row.Columns[4].GetString_()
	po.SubOrder_Id = row.Columns[5].GetString_()
	po.Order_Status = row.Columns[6].GetString_()
	po.Asset_ID = row.Columns[7].GetString_()
	return po
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	err = stub.DelState("role_OEM")
	if err != nil {
		return nil, fmt.Errorf("remove operation failed. Error updating state: %s", err)
	}

	err = stub.DelState("role_tier_1")
	if err != nil {
		return nil, fmt.Errorf("remove operation failed. Error updating state: %s", err)
	}
	err = stub.DelState("role_first_tier_2")
	if err != nil {
		return nil, fmt.Errorf("remove operation failed. Error updating state: %s", err)
	}
	err = stub.DelState("role_second_tier_2")
	if err != nil {
		return nil, fmt.Errorf("remove operation failed. Error updating state: %s", err)
	}

	err = stub.CreateTable("PurchaseOrder", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Order_Id", shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Order_Desc", shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Order_Quantity", shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Assigned_To_Id", shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Created_By_Id", shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "SubOrder_Id", shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Order_Status", shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "Asset_ID", shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating PurchaseOrder Table.")
	}
	fmt.Println("IN init table created successfully %s", err)
	return nil, nil

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "createOrder" {
		fmt.Println("IN invoke functions ==========")

		return t.createOrder(stub, args)
	}
	if function == "updateOrderStatus" {

		return updateOrderStatus(stub, args)
	}
	return nil, errors.New("Received unknown function invocation: " + function)

}
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "fetchAllOrders" {
		fmt.Println("IN OUERY ==========================")

		return t.fetchAllOrders(stub, args)
	}

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) createOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("In create order")

	col1Val := args[0]
	col2Val := args[1]
	col3Val := args[2]
	col4Val := args[3]
	col5Val := args[4]
	col6Val := args[5]
	col7Val := args[6]
	col8Val := args[7]

	ok, err := stub.InsertRow("PurchaseOrder", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: col1Val}},
			&shim.Column{Value: &shim.Column_String_{String_: col2Val}},
			&shim.Column{Value: &shim.Column_String_{String_: col3Val}},
			&shim.Column{Value: &shim.Column_String_{String_: col4Val}},
			&shim.Column{Value: &shim.Column_String_{String_: col5Val}},
			&shim.Column{Value: &shim.Column_String_{String_: col6Val}},
			&shim.Column{Value: &shim.Column_String_{String_: col7Val}},
			&shim.Column{Value: &shim.Column_String_{String_: col8Val}},
		}})

	if err != nil {
		return nil, err
	}
	if !ok && err == nil {
		return nil, errors.New("Row already exists.")
	}
	fmt.Println("Value of Columns ====" + col1Val + ", " + col2Val + " ," + col3Val + "," + col4Val + "," + col5Val + "," + col6Val + "," + col7Val + "," + col8Val)
	fmt.Println("After row Inserted==========%s", ok)
	return []byte("success"), errors.New("Received unknown function invocation: ")
}

func (t *SimpleChaincode) fetchAllOrders(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("IN FETCH ALL ORDERS============================")
	//all  id with overall status(irrespective of the role)
	var columns []shim.Column
	col1Val := args[0]
	col1 := shim.Column{Value: &shim.Column_String_{String_: col1Val}}
	columns = append(columns, col1)
	fmt.Println("Befor get Rows")
	row, err := stub.GetRow("PurchaseOrder", columns)
	if err != nil {
		fmt.Println("hi")
		return nil, fmt.Errorf("Failed to retrieve row")
	}
	//orderArrary := []*PO_tier1{}
	var po = PO_tier1{}
	if len(row.Columns) < 1 {
		fmt.Println("hello")
		return []byte("no row found with id:" + args[0]), fmt.Errorf("Failed to retrieve row")
	}
	//for row := range rows {
	//po = new(PO_tier1)
	po.Order_Id = row.Columns[0].GetString_()
	fmt.Println("PO ID====" + row.Columns[0].GetString_())
	po.Order_Desc = row.Columns[1].GetString_()
	po.Order_Quantity = row.Columns[2].GetString_()
	po.Assigned_To_Id = row.Columns[3].GetString_()
	po.Created_By_Id = row.Columns[4].GetString_()
	po.SubOrder_Id = row.Columns[5].GetString_()
	fmt.Println("PO SUBID====" + row.Columns[5].GetString_())
	po.Order_Status = row.Columns[6].GetString_()
	po.Asset_ID = row.Columns[7].GetString_()

	//		orderArrary = append(orderArrary, po)
	//}

	jsonRows, _ := json.Marshal(po)
	fmt.Println("Printing rows==========================" + string(jsonRows))
	return jsonRows, nil

}

func updateOrderStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	return nil, nil
}
