package parser

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	jsonstruct "study/internal/parser/models/json_struct"
)

func UrlParse(ID int, off int) *http.Response {
	reqBody, err := json.Marshal(jsonstruct.Request{
		CategoryID: ID,
		Filters: jsonstruct.Filters{
			Checkbox:      []interface{}{},
			Multicheckbox: []interface{}{},
			Range:         []interface{}{},
		},
		Sort: jsonstruct.Sort{
			Type:  "popular",
			Order: "desc",
		},
		Limit:  40,
		Offset: off,
	})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "https://lenta.com/api-gateway/v1/catalog/items", bytes.NewBuffer(reqBody))
	req.Header.Set("deviceid", "05d46805-11a1-846e-578c-84c63e768547")
	req.Header.Set("x-platform", "omniweb")
	req.Header.Set("x-retail-brand", "lo")
	req.Header.Set("Cookie", "GrowthBook_user_id=63d3ae61-1309-f9b6-0427-1fb992ce08f5; App_Cache_MPK=mp300-b1de0bac2c257f3257bf5ef2eea4ecbc; App_Cache_CitySlug=moscow; UserSessionId=6a0f1434-9ff2-9cf1-8f6e-4afba77a85eb; Utk_SessionToken=08D18C20128C0DC83C8E003061270FF8; App_Cache_City=%7B%22centerLat%22%3A%2255.75322000%22%2C%22centerLng%22%3A%2237.62255200%22%2C%22id%22%3A1%2C%22isDefault%22%3Atrue%2C%22mainDomain%22%3Afalse%2C%22name%22%3A%22%D0%9C%D0%BE%D1%81%D0%BA%D0%B2%D0%B0%20%D0%B8%20%D0%9C%D0%9E%22%2C%22slug%22%3A%22moscow%22%7D; oxxfgh=d7980825-0d4a-4264-95a9-b9b69133e98e%230%235184000000%235000%231800000%2344965; uwyii=60c01ce0-35da-66c5-d4a3-830c05c0b88b; iap.uid=4dc9f8412eab4e358bc6067b7b2f1993; _ym_uid=1772479011718506031; _ym_d=1772479011; flocktory-uuid=d937906e-bddf-44dd-a66d-858dc51fc1ac-3; tmr_lvid=1e5002cb9cffea91255b0ad9ea76e19b; tmr_lvidTS=1772479010906; agree_with_cookie=true; Utk_MrkGrpTkn=0D68EA9CE36434E0C55AF63C116E28A0; _ym_isad=1; domain_sid=NtKBFUm2bt_POQBfWBfrF%3A1773668488313; Utk_DvcGuid=05d46805-11a1-846e-578c-84c63e768547; User_Agent=Mozilla%2F5.0%20(Windows%20NT%2010.0%3B%20Win64%3B%20x64)%20AppleWebKit%2F537.36%20(KHTML%2C%20like%20Gecko)%20Chrome%2F145.0.0.0%20Safari%2F537.36; Is_Search_Bot=false; App_Cache_MissionAddressMode=%7B%22t%22%3A%22pickup%22%2C%22ids%22%3Afalse%2C%22ma%22%3A%7B%22i%22%3A104%2C%22a%22%3A%220614%22%2C%22t%22%3A%22%D0%A2%D0%9A614%22%2C%22af%22%3A%22%D0%9C%D0%BE%D1%81%D0%BA%D0%B2%D0%B0%2C%20%D0%93%D1%83%D1%80%D1%8C%D1%8F%D0%BD%D0%BE%D0%B2%D0%B0%20%D1%83%D0%BB.%2C%202%D0%90%22%2C%22ri%22%3A1%2C%22mt%22%3A%22SM%22%2C%22s%22%3Afalse%7D%7D; CatalogSelectedSorting=popular; GrowthBook_experiments=experiment_web_aa_test_202601.1%2Cexperiment_web_new_navigation_web_mobile_2026_v3.0%2Cexperiment_search_collections_ranking_web_2026_03_13.0%2Cexperiment_web_leave_order_at_door.0%2Cexperiment_web_aa_2026_01_v1.0%2Cexperiment_web_ds_cat_diversity_v2.1; GrowthBook_Cookie_Experiments=exp_newui_chips.test%2C%20exp_web_chips_online.default%2C%20exp_chips_online.default%2C%20exp_big_card.Default%2C%20exp_product_page_by_blocks.test%2C%20exp_without_a_doorbell.default%2C%20exp_without_a_doorbell_new.default%2C%20exp_search_photo_positions.default%2C%20exp_new_navigation_web.test%2C%20exp_web_mobile_tabbar.default%2C%20exp_web_personal_promo_delivery_chips.test%2C%20exp_personal_promo_delivery_chips.default%2C%20exp_new_navigation_web_search.control%2C%20exp_new_navigation_web_actions.control%2C%20exp_leave_order_at_door.control%2C%20exp_leave_order_at_door_new.control%2C%20exp_unpin_tabbar.default%2C%20exp_unpin_tabbar_v2.default; qrator_jsid=1773665405.300.5PQCxWBmMiW3qlpK-2h2f11qciisiv5of48a4f7nli8e39rvs; spses.d58d=*; _ym_visorc=b; uwyiert=283d4a04-3246-dc88-854c-ae696a08d485; spid.d58d=694fa5a5-1f2d-460c-9f42-7d4e6549a9c4.1772479011.9.1773701257.1773695518.9f4bfae3-3fe2-46d8-8649-ae59028b79bd.22d16b9e-c943-4484-b90c-56e6f4799623.411fb377-7c0b-43a6-8e4c-421d755eaca2.1773701247159.27; tmr_detect=1%7C1773701256651")
	req.Header.Set("sessiontoken", "08D18C20128C0DC83C8E003061270FF8")

	c := &http.Client{}
	res, _ := c.Do(req)

	return res
}

func byteConv(res *http.Response) []byte {
	body, _ := io.ReadAll(res.Body)

	return body
}

func CategoryParse() map[string]int {
	var body jsonstruct.Response
	categories := make(map[string]int, 0)

	json.Unmarshal(byteConv(UrlParse(1893, 0)), &body)

	for _, v := range body.Categories {
		categories[v.Name] = v.ID
	}

	return categories
}

func ProductCardParse(body jsonstruct.JsonProduct, p *[]jsonstruct.Product) {
	for _, v := range body.Items {
		*p = append(*p, jsonstruct.Product{
			Name:     v.Name,
			OldPrice: v.Prices.Price / 100,
			NewPrice: v.Prices.PriceRegular / 100,
		})
	}
}

func Pagination(ID int) []jsonstruct.Product {
	var body jsonstruct.JsonProduct
	products := make([]jsonstruct.Product, 0)

	for i := 0; i > -1; i += 40 {
		if res := UrlParse(ID, i); res.StatusCode != 200 {
			return products
		} else {
			json.Unmarshal(byteConv(res), &body)
			ProductCardParse(body, &products)
		}
	}

	return nil
}
