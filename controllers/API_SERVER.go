package APIServer
import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

)



var (
	idNotPresented        = errors.New("ID du compte non présent dans les parametres de la requete")
	timeRangeNotPresented = errors.New("le nombre de jours avant que la requete ne transfère les statistiques n'est pas present dans les parametres de la requete")
	limitNotPresented     = errors.New("la limite de la requete n'est pas presente dans les parametres de la requete")
)

// API_SERVER contient les données nécessaires pour éxécuter le serveur API
type APIServer struct {
	config *Config    //configuration
	logger *logger    //enregistreur
	router *http.ServeMux   //routeur
	store  store.Store
}

// Creation de la nouvelle instance de la structure de API
func New(config *Config) *APIServer {
	s := &APIServer{
		config: config,
		logger: NewLogger(),
		router: http.NewServeMux(),
	}
	s.configureLogger()
	s.configureRouter()
	return s
}

func (s *APIServer) setStore(store store.Store) {
	s.store = store
}

// Demarrage de l'execution du serveur de base de données de API
func (s *APIServer) Start() error {
	store, err := sqlstore.New(s.config.DbPath, s.config.QueryTimeout)
	if err != nil {
		return nil
	}
	defer store.Close()
	s.setStore(store)
	s.logger.Info("Demarrage du serveur api")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() {
	s.logger.SetLevel(s.config.LogLevel)
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/health", s.handleHealth())
	s.router.HandleFunc("/api/v1/accounts", s.handleAccounts())
	s.router.HandleFunc("/api/v1/transfer-money", s.handleTransferMoney())
	s.router.HandleFunc("/api/v1/transactions", s.handleTransactions())
}

func (s *APIServer) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func (s *APIServer) handleError(err error, statusCode int, w http.ResponseWriter, r *http.Request) {
	s.logger.Error(fmt.Sprintf("Methode: %s; erreur: %s", r.URL.Path, err.Error()))
	w.WriteHeader(statusCode)
}

func parseIntQueryParams(r *http.Request, paramNames ...string) (map[string]int64, error) {
	params := r.URL.Query()
	m := make(map[string]int64)
	for _, param := range paramNames {
		val := params.Get(param)
		conv, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		m[param] = conv
	}
	return m, nil
}

func (s *APIServer) handleAccounts() http.HandlerFunc {  //  la fonction qui assure la gestion des comptes
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":    //POST est un verbe comme GET, DELETE, PATCH, INPUT, il est utilisé pour la transmission des données
			w.Header().Set("Content-type", "application/json")
			var acc AccountJsonView
			err := json.NewDecoder(r.Body).Decode(&acc)
			if err != nil {
				s.handleError(err, http.StatusBadRequest, w, r)
				return
			}
			accModel, err := s.store.InsertAccount(acc.Balance)
			if err != nil {
				s.handleError(err, http.StatusInternalServerError, w, r)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(AccountIDJsonView{ID: accModel.AccountID})
		case "DELETE"://DELETE est un verbe comme GET, PATCH, INPUT, il est utilisé pour la suppression des données
			valMap, err := parseIntQueryParams(r, "account_id")
			if err != nil {
				s.handleError(err, http.StatusBadRequest, w, r)
				return
			}
			err = s.store.DeleteAccount(valMap["account_id"])
			if err != nil {
				s.handleError(err, http.StatusInternalServerError, w, r)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		case "GET"://GET est un verbe comme  DELETE, PATCH, INPUT, il est utilisé pour la recuperation des données
			w.Header().Set("Content-type", "application/json")
			valMap, err := parseIntQueryParams(r, "account_id")
			if err != nil {
				s.handleError(err, http.StatusBadRequest, w, r)
				return
			}
			accModel, err := s.store.GetAccount(valMap["account_id"])
			if err != nil {
				s.handleError(err, http.StatusInternalServerError, w, r)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(AccountJsonView{
				AccountID: accModel.AccountID,
				Balance:   accModel.Balance,
			})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		}
	}
}

func (s *APIServer) handleTransferMoney() http.HandlerFunc {  //  la fonction qui assure la gestion des transferts d'argent
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {  //
		case "POST":
			var tr TransactionJsonView
			err := json.NewDecoder(r.Body).Decode(&tr)
			if err != nil {
				s.handleError(err, http.StatusBadRequest, w, r) //gestion des erreurs
			}
			err = s.store.TransferMoney(
				tr.ToAccountID,
				tr.FromAccountID,
				tr.Amount,
			)
			if err != nil {
				s.handleError(err, http.StatusInternalServerError, w, r)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		}
	}
}

func (s *APIServer) handleTransactions() http.HandlerFunc {  // la fonction qui assure la gestion des transactions
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-type", "application/json")
			valMap, err := parseIntQueryParams(r, "account_id", "n_last_days", "limit")
			if err != nil {
				s.handleError(err, http.StatusBadRequest, w, r)  //gestion des erreurs
				return
			}

			transactions, err := s.store.GetTransactionsHistory(  //ici on récupère l'historique des transactions
				valMap["account_id"],  //on recupère id du compte
				valMap["n_last_days"],  //on récupère la jour de la requete
				valMap["limit"],      //on recupère la limit de depot d'argent
			)
			if err != nil {
				s.handleError(err, http.StatusInternalServerError, w, r)
				return
			}
			transactionsJson := make([]TransactionJsonView, len(transactions))
			for i, tr := range transactions {  //ici on range dans un tableau
				transactionsJson[i] = TransactionJsonView{
					Timestamp:     tr.Timestamp,
					FromAccountID: tr.FromAccountID,
					ToAccountID:   tr.ToAccountID,
					Amount:        tr.Amount,
				}
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(transactionsJson)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		}
	}
}
