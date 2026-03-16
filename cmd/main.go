package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	CategoryID int     `json:"categoryId"`
	Filters    Filters `json:"filters"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
	Sort       Sort    `json:"sort"`
}
type Filters struct {
	Checkbox      []interface{} `json:"checkbox"`
	Multicheckbox []interface{} `json:"multicheckbox"`
	Range         []interface{} `json:"range"`
}
type Sort struct {
	Type  string `json:"type"`
	Order string `json:"order"`
}

func main() {
	// открыть https://lenta.com/catalog/osobenno-vygodno-1893/
	reqBody, err := json.Marshal(Request{
		CategoryID: 1893,
		Filters: Filters{
			Checkbox:      []interface{}{},
			Multicheckbox: []interface{}{},
			Range:         []interface{}{},
		},
		Sort: Sort{
			Type:  "popular",
			Order: "desc",
		},
		Limit:  40,
		Offset: 0,
	})
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", "https://lenta.com/api-gateway/v1/catalog/items", bytes.NewBuffer(reqBody))

	req.Header.Set("deviceid", "fb11183f-4d0e-006f-b043-6378283eac56")
	req.Header.Set("x-platform", "omniweb")
	req.Header.Set("x-retail-brand", "Io")
	req.Header.Set("Cookie", "GrowthBook_user_id=63d3ae61-1309-f9b6-0427-1fb992ce08f5; App_Cache_MPK=mp300-b1de0bac2c257f3257bf5ef2eea4ecbc; App_Cache_CitySlug=moscow; UserSessionId=6a0f1434-9ff2-9cf1-8f6e-4afba77a85eb; Utk_SessionToken=08D18C20128C0DC83C8E003061270FF8; App_Cache_City=%7B%22centerLat%22%3A%2255.75322000%22%2C%22centerLng%22%3A%2237.62255200%22%2C%22id%22%3A1%2C%22isDefault%22%3Atrue%2C%22mainDomain%22%3Afalse%2C%22name%22%3A%22%D0%9C%D0%BE%D1%81%D0%BA%D0%B2%D0%B0%20%D0%B8%20%D0%9C%D0%9E%22%2C%22slug%22%3A%22moscow%22%7D; App_Cache_MissionAddressMode=%7B%22t%22%3A%22pickup%22%2C%22ids%22%3Atrue%2C%22ma%22%3A%7B%22i%22%3A3149%2C%22a%22%3A%220124%22%2C%22t%22%3A%22%D0%A2%D0%9A124%22%2C%22af%22%3A%22%D0%9C%D0%BE%D1%81%D0%BA%D0%B2%D0%B0%2C%207-%D1%8F%20%D0%9A%D0%BE%D0%B6%D1%83%D1%85%D0%BE%D0%B2%D1%81%D0%BA%D0%B0%D1%8F%20%D1%83%D0%BB%D0%B8%D1%86%D0%B0%2C%209%20(%D0%A2%D0%A0%D0%A6%20%D0%9C%D0%BE%D0%B7%D0%B0%D0%B8%D0%BA%D0%B0)%22%2C%22ri%22%3A1%2C%22mt%22%3A%22HM%22%2C%22s%22%3Afalse%7D%7D; oxxfgh=d7980825-0d4a-4264-95a9-b9b69133e98e%230%235184000000%235000%231800000%2344965; uwyii=60c01ce0-35da-66c5-d4a3-830c05c0b88b; iap.uid=4dc9f8412eab4e358bc6067b7b2f1993; _ym_uid=1772479011718506031; _ym_d=1772479011; flocktory-uuid=d937906e-bddf-44dd-a66d-858dc51fc1ac-3; tmr_lvid=1e5002cb9cffea91255b0ad9ea76e19b; tmr_lvidTS=1772479010906; domain_sid=NtKBFUm2bt_POQBfWBfrF%3A1773501396986; agree_with_cookie=true; CatalogSelectedSorting=popular; spid.d58d=694fa5a5-1f2d-460c-9f42-7d4e6549a9c4.1772479011.3.1773507710.1773503010.4bccdbc7-428d-4731-8793-63e521c7624e.43a2b48b-38c1-401c-a33a-81957f870322.8c0fcbb5-8f69-48e3-ba1c-c51c835f4532.1773506674397.160; qrator_jsr=1773665405.300.5PQCxWBmMiW3qlpK-1d2ap0da01v2effnr2ajfvnndvpo06eg-00; qrator_jsid=1773665405.300.5PQCxWBmMiW3qlpK-ch3rcl5imkpkk9a2jd1bubkdvnp2hti5; Utk_DvcGuid=fb11183f-4d0e-006f-b043-6378283eac56; User_Agent=Mozilla%2F5.0%20(Windows%20NT%2010.0%3B%20Win64%3B%20x64)%20AppleWebKit%2F537.36%20(KHTML%2C%20like%20Gecko)%20Chrome%2F145.0.0.0%20Safari%2F537.36; Is_Search_Bot=false; Utk_MrkGrpTkn=0D68EA9CE36434E0C55AF63C116E28A0; uwyiert=283d4a04-3246-dc88-854c-ae696a08d485; GrowthBook_experiments=experiment_web_aa_test_202601.1%2Cexperiment_web_new_navigation_web_mobile_2026_v3.0%2Cexperiment_search_collections_ranking_web_2026_03_13.0%2Cexperiment_web_leave_order_at_door.0%2Cexperiment_web_aa_2026_01_v1.0%2Cexperiment_web_ds_cat_diversity_v2.1; GrowthBook_Cookie_Experiments=exp_newui_chips.test%2C%20exp_web_chips_online.default%2C%20exp_chips_online.default%2C%20exp_big_card.Default%2C%20exp_product_page_by_blocks.test%2C%20exp_search_photo_positions.default%2C%20exp_new_navigation_web.test%2C%20exp_web_mobile_tabbar.default%2C%20exp_web_personal_promo_delivery_chips.test%2C%20exp_personal_promo_delivery_chips.default%2C%20exp_new_navigation_web_search.control%2C%20exp_new_navigation_web_actions.control%2C%20exp_unpin_tabbar.default%2C%20exp_unpin_tabbar_v2.default")
	req.Header.Set("sessiontoken", "08D18C20128C0DC83C8E003061270FF8")

	c := &http.Client{}
	res, _ := c.Do(req)
	body, _ := io.ReadAll(res.Body)
	fmt.Println(len(body))

	var response map[string]interface{}

	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	p, _ := json.MarshalIndent(response, "", " ")
	fmt.Println(string(p))
}
